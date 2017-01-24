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
	// Templates
	s.addTemplatesRoute()

	// Routes
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

	// Static assets
	if len(s.StaticDirectory) > 0 {
		s.FileHandler = http.FileServer(http.Dir(s.StaticDirectory))
		staticDirectory := "/" + s.StaticDirectory + "/"
		s.HttpServerMux.Handle(staticDirectory, http.StripPrefix(staticDirectory, s.FileHandler))
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
