package web

import (
	"net/http"
	"fmt"
	"strconv"
)

type Server struct {
	Port          int
	TemplateDir   string
	StaticDir     string
	Routes        []Route
	Sessions      bool
	HttpServerMux *http.ServeMux
	FileHandler   http.Handler
}

func (s *Server) Run() {
	s.addTemplatesRoute()
	s.addStaticRoute()
	s.setupHandlers()
	s.startServer()
}

func (s *Server) addTemplatesRoute() {
	if len(s.TemplateDir) > 0 {
		s.Routes = append(s.Routes, Route{
			Pattern: "/",
			Handler: func(r *Request) {
				renderer, err := GetRenderer(s.TemplateDir)
				check(err)

				r.HttpResponseWriter.Header().Set("Content-Type", "text/html")

				filename := GetFilenameFromRequest(r.HttpRequest)

				if len(filename) == 0 {
					filename = "index"
				}

				renderer.Render([]string{
					filename + ".html",
					"404.html",
				}, r.HttpResponseWriter, r)
			},
		})
	}
}

func (s *Server) addStaticRoute() {
	if len(s.StaticDir) > 0 {
		staticDirectory := s.StaticDir + "/"
		s.Routes = append(s.Routes, Route{
			Pattern: staticDirectory,
			Handler: func(r *Request) {
				s.FileHandler = http.FileServer(http.Dir(s.StaticDir))
				handler := http.StripPrefix(staticDirectory, s.FileHandler)
				handler.ServeHTTP(r.HttpResponseWriter, &r.HttpRequest)
			},
		})
	}
}

func (s *Server) setupHandlers() {
	s.HttpServerMux = http.NewServeMux()
	for _, routeTemp := range s.Routes {
		route := routeTemp
		fmt.Printf("Setting pattern: %s\n", route.Pattern)
		s.HttpServerMux.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			request := Request{
				HttpResponseWriter: w,
				HttpRequest: *r,
			}
			if s.Sessions {
				request.InitSession()
			}
			if route.CsrfProtect && ! request.IsCsrfPresentAndValid() {
				request.SetResponseCode(http.StatusForbidden)
			} else {
				route.Handler(&request)
			}
			fmt.Printf("Handled request: %#v\n", r.URL.Path)
		})
	}
}

func (s *Server) startServer() {
	err := http.ListenAndServe(":" + strconv.Itoa(s.Port), s.HttpServerMux)
	check(err)
}
