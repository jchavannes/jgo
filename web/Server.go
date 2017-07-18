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
	Port              int
	Routers           []Router
	SessionKey        string
	StaticFilesDir    string
	StrictSlash       bool
	TemplatesDir      string
	UseAutoRender     bool
	UseSessions       bool
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
	for _, routerTemp := range s.Routers {
		router := routerTemp
		prefix := router.GetPatternPrefix()
		s.router.PathPrefix(prefix).Handler(Handler{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				response := getResponse(w, r, s)
				defer response.LogComplete()
				if router.PreHandler != nil {
					router.PreHandler(&response)
				}

				if prefix == "/" {
					if len(s.StaticFilesDir) > 0 {
						allowedFileTypes := GetDefaultAllowedFileExtensions()
						if len(s.AllowedExtensions) > 0 {
							allowedFileTypes = s.AllowedExtensions
						}
						for _, fileType := range allowedFileTypes {
							if strings.HasSuffix(response.Request.HttpRequest.URL.Path, "." + fileType) || fileType == "*" {
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
				}

				if router.NotFoundHandler != nil {
					router.NotFoundHandler(&response)
				}
				response.SetResponseCode(http.StatusNotFound)
			},
		})
	}
}

func (s *Server) setupHandlers() {
	s.router = mux.NewRouter()
	s.router.StrictSlash(s.StrictSlash)
	for _, routerTemp := range s.Routers {
		router := routerTemp
		for _, routeTemp := range router.Routes {
			route := routeTemp
			name := ""
			if len(route.Name) > 0 {
				name = " (" + route.Name + ")"
			}
			pattern := router.GetRoutePattern(route)
			fmt.Printf("Adding pattern to router: %s%s\n", pattern, name)
			s.router.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
				response := getResponse(w, r, s)
				if router.PreHandler != nil {
					router.PreHandler(&response)
				}
				defer response.LogComplete()
				if route.CsrfProtect && ! response.IsValidCsrf() {
					response.SetResponseCode(http.StatusForbidden)
				} else {
					route.Handler(&response)
				}
			})
		}
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
	return response
}

func (s *Server) startServer() error {
	srv := &http.Server{
		Handler: s.router,
		Addr: ":" + strconv.Itoa(s.Port),
	}
	fmt.Printf("Starting server on port %d...\n", s.Port)
	return srv.ListenAndServe()
}
