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
	HttpServerMux     *http.ServeMux
	FileHandler       http.Handler
}

func (s *Server) Run() {
	s.Routes = append(s.Routes, Route{
		Pattern: "/",
		Handler: func(r *Request) {
			renderer, err := GetRenderer(s.TemplateDirectory)
			check(err)

			writer := *r.HttpResponseWriter
			writer.Header().Set("Content-Type", "text/html")

			filename := GetFilenameFromRequest(r.HttpRequest)

			if len(filename) == 0 {
				filename = "index"
			}

			renderer.Render([]string{
				filename + ".html",
				"404.html",
			}, *r.HttpResponseWriter, nil)
		},
	})

	// Requests
	s.HttpServerMux = http.NewServeMux()
	for _, route := range s.Routes {
		fmt.Printf("Setting pattern: %s\n", route.Pattern)
		handler := route.Handler
		s.HttpServerMux.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			handler(&Request{
				HttpResponseWriter: &w,
				HttpRequest: r,
			})
			fmt.Printf("Handled request: %#v\n", r.URL.Path)
		})
	}

	// Static assets
	s.FileHandler = http.FileServer(http.Dir(s.StaticDirectory))
	staticDirectory := "/" + s.StaticDirectory + "/"
	s.HttpServerMux.Handle(staticDirectory, http.StripPrefix(staticDirectory, s.FileHandler))

	err := http.ListenAndServe(":" + strconv.Itoa(s.Port), s.HttpServerMux)
	check(err)
}
