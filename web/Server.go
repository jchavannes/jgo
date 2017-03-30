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
	DisableAutoRender bool
	Port              int
	PreHandler        func(*Response)
	Routes            []Route
	SessionKey        string
	StaticFilesDir    string
	TemplatesDir      string
	UseSessions       bool
	router            *mux.Router
}

var (
	AllFileExtensions = []string{
		"*",
	}

	DefaultAllowedFileExtensions = []string{
		"js",
		"css",
		"jpg",
		"png",
		"ico",
	}
)

func (s *Server) Run() {
	s.addTemplatesRoute()
	s.setupHandlers()
	s.startServer()
}

func (s *Server) addTemplatesRoute() {
	if len(s.TemplatesDir) > 0 && ! s.DisableAutoRender {
		s.Routes = append(s.Routes, Route{
			Pattern: "/",
			Name: "Automatic renderer - templates directory: " + s.TemplatesDir,
			Handler: func(r *Response) {
				templateName := r.Request.GetPotentialFilename()

				if len(templateName) == 0 {
					templateName = "index"
				}

				r.RenderTemplate(templateName)
			},
		})
	}
}

func (s *Server) setupHandlers() {
	s.router = mux.NewRouter()
	for _, routeTemp := range s.Routes {
		route := routeTemp
		name := ""
		if len(route.Name) > 0 {
			name = " (" + route.Name + ")"
		}
		fmt.Printf("Adding pattern to router: %s%s\n", route.Pattern, name)
		s.router.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
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
			}
			if s.PreHandler != nil {
				s.PreHandler(&response)
			}
			response.Helper["CsrfToken"] = response.Session.GetCsrfToken()
			if route.CsrfProtect && ! response.IsValidCsrf() {
				response.SetResponseCode(http.StatusForbidden)
			} else {
				route.Handler(&response)
			}
			fmt.Printf("Handled request: %#v\n", r.URL.Path)
		})
	}
	if len(s.StaticFilesDir) > 0 {
		fmt.Printf("Adding static file handler: %s\n", s.StaticFilesDir)
		handler := Handler{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				allowedFileTypes := DefaultAllowedFileExtensions
				if len(s.AllowedExtensions) > 0 {
					allowedFileTypes = s.AllowedExtensions
				}
				for _, fileType := range allowedFileTypes {
					if strings.HasSuffix(r.URL.Path, "." + fileType) || fileType == "*" {
						http.FileServer(http.Dir(s.StaticFilesDir)).ServeHTTP(w, r)
						return
					}
				}
			},
		}
		s.router.PathPrefix("/").Handler(handler)
	}
}

func (s *Server) startServer() {
	srv := &http.Server{
		Handler: s.router,
		Addr: ":" + strconv.Itoa(s.Port),
	}
	fmt.Printf("Starting server on port %d...\n", s.Port)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
