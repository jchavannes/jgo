package web

import (
	"net/http"
	"fmt"
	"strconv"
)

type Server struct {
	Port          int
	Template      Template
	Static        Static
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
	if len(s.Template.Directory) > 0 {
		s.Routes = append(s.Routes, Route{
			Pattern: s.Template.Pattern,
			Handler: func(r *Request) {
				renderer, err := GetRenderer(s.Template.Directory)
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
	if len(s.Static.Directory) > 0 {
		s.Routes = append(s.Routes, Route{
			Pattern: s.Static.Pattern,
			Handler: func(r *Request) {
				s.FileHandler = http.FileServer(http.Dir(s.Static.Directory))
				handler := http.StripPrefix("/" + s.Static.Directory + "/", s.FileHandler)
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
