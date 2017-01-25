package web

import (
	"net/http"
	"fmt"
	"strconv"
)

type Server struct {
	Port              int
	TemplateDirectory string
	StaticDirectory   string
	Routes            []Route
	EnableSessions    bool
	HttpServerMux     *http.ServeMux
	FileHandler       http.Handler
}

func (s *Server) Run() {
	s.addTemplatesRoute()
	s.addStaticRoute()

	s.HttpServerMux = http.NewServeMux()
	for _, routeTemp := range s.Routes {
		route := routeTemp
		fmt.Printf("Setting pattern: %s\n", route.Pattern)
		s.HttpServerMux.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			request := Request{
				HttpResponseWriter: w,
				HttpRequest: *r,
			}
			if s.EnableSessions {
				request.InitSession()
			}
			if route.CsrfProtect {

			}
			route.Handler(&request)
			fmt.Printf("Handled request: %#v\n", r.URL.Path)
		})
	}

	err := http.ListenAndServe(":" + strconv.Itoa(s.Port), s.HttpServerMux)
	check(err)
}

func (s *Server) addTemplatesRoute() {
	if len(s.TemplateDirectory) > 0 {
		s.Routes = append(s.Routes, Route{
			Pattern: "/",
			Handler: func(r *Request) {
				renderer, err := GetRenderer(s.TemplateDirectory)
				check(err)

				r.HttpResponseWriter.Header().Set("Content-Type", "text/html")

				filename := GetFilenameFromRequest(r.HttpRequest)

				if len(filename) == 0 {
					filename = "index"
				}

				renderer.Render([]string{
					filename + ".html",
					"404.html",
				}, r.HttpResponseWriter, nil)
			},
		})
	}
}
func (s *Server) addStaticRoute() {
	if len(s.StaticDirectory) > 0 {
		staticDirectory := "/" + s.StaticDirectory + "/"
		s.Routes = append(s.Routes, Route{
			Pattern: staticDirectory,
			Handler: func(r *Request) {
				s.FileHandler = http.FileServer(http.Dir(s.StaticDirectory))
				handler := http.StripPrefix(staticDirectory, s.FileHandler)
				handler.ServeHTTP(r.HttpResponseWriter, &r.HttpRequest)
			},
		})
	}
}
