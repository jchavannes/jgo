package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"strconv"
	"strings"
)

type Server struct {
	AllowedExtensions []string
	NotFoundHandler   func(*Response)
	Port              int
	PreHandler        func(*Response)
	Routes            []Route
	SessionKey        string
	StaticFilesDir    string
	StrictSlash       bool
	TemplatesDir      string
	UseAutoRender     bool
	UseSessions       bool
	InsecureCookie    bool
	router            *mux.Router
}

// Default extensions allowed for static files.
func GetDefaultAllowedFileExtensions() []string {
	return []string{
		"js",
		"css",
		"jpg",
		"png",
		"ico",
		"gif",
	}
}

func (s *Server) Run() error {
	s.setupHandlers()
	s.addCatchAllRoute()
	return s.startServer()
}

func (s *Server) addCatchAllRoute() {
	if len(s.StaticFilesDir) > 0 {
		fmt.Printf("Static files directory: %s\n", s.StaticFilesDir)
	}
	if len(s.TemplatesDir) > 0 {
		fmt.Printf("Templates directory: %s\n", s.TemplatesDir)
	}
	s.router.PathPrefix("/").Handler(Handler{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response := getResponse(w, r, s)
			defer response.LogComplete()

			if len(s.StaticFilesDir) > 0 {
				allowedFileTypes := GetDefaultAllowedFileExtensions()
				if len(s.AllowedExtensions) > 0 {
					allowedFileTypes = s.AllowedExtensions
				}
				for _, fileType := range allowedFileTypes {
					if strings.HasSuffix(response.Request.HttpRequest.URL.Path, "."+fileType) || fileType == "*" {
						http.FileServer(http.Dir(s.StaticFilesDir)).ServeHTTP(response.Writer, &response.Request.HttpRequest)
						return
					}
				}
			}

			if len(s.TemplatesDir) > 0 && s.UseAutoRender {
				templateName := response.Request.GetPotentialFilename()

				if len(templateName) == 0 {
					templateName = "index"
				}

				err := response.RenderTemplate(templateName)
				if err == nil {
					return
				}
			}

			if s.NotFoundHandler != nil {
				s.NotFoundHandler(&response)
			} else {
				response.SetResponseCode(http.StatusNotFound)
			}
		},
	})
}

func (s *Server) setupHandlers() {
	s.router = mux.NewRouter()
	s.router.StrictSlash(s.StrictSlash)
	for _, routeTemp := range s.Routes {
		route := routeTemp
		name := ""
		if len(route.Name) > 0 {
			name = " (" + route.Name + ")"
		}
		fmt.Printf("Adding pattern to router: %s%s\n", route.Pattern, name)
		s.router.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			response := getResponse(w, r, s)
			defer response.LogComplete()
			if route.CsrfProtect && ! response.IsValidCsrf() {
				response.SetResponseCode(http.StatusForbidden)
			} else {
				route.Handler(&response)
			}
		})
	}
}

func getResponse(w http.ResponseWriter, r *http.Request, s *Server) Response {
	response := Response{
		Helper: make(map[string]interface{}),
		Writer: w,
		Request: Request{
			HttpRequest: *r,
		},
		Server: s,
	}
	if s.UseSessions {
		response.InitSession()
		response.Helper["CsrfToken"] = response.Session.GetCsrfToken()
	}
	if s.PreHandler != nil {
		s.PreHandler(&response)
	}
	return response
}

func (s *Server) startServer() error {
	srv := &http.Server{
		Handler: s.router,
		Addr:    ":" + strconv.Itoa(s.Port),
	}
	fmt.Printf("Starting server on port %d...\n", s.Port)
	return srv.ListenAndServe()
}
