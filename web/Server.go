package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"strconv"
)

type Server struct {
	InitResponse func(*Response)
	Port         int
	Router       *mux.Router
	Routes       []Route
	SessionKey   string
	Sessions     bool
	StaticDir    string
	TemplateDir  string
}

func (s *Server) Run() {
	s.addTemplatesRoute()
	s.setupHandlers()
	s.startServer()
}

func (s *Server) addTemplatesRoute() {
	if len(s.TemplateDir) > 0 {
		s.Routes = append(s.Routes, Route{
			Pattern: "/",
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
	s.Router = mux.NewRouter()
	for _, routeTemp := range s.Routes {
		route := routeTemp
		fmt.Printf("Setting pattern: %s\n", route.Pattern)
		s.Router.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			response := Response{
				Helper: make(map[string]interface{}),
				Writer: w,
				Request: Request{
					HttpRequest: *r,
				},
				TemplateDir: s.TemplateDir,
			}
			if s.Sessions {
				response.InitSession(s.SessionKey)
			}
			s.InitResponse(&response)
			response.Helper["CsrfToken"] = response.Session.GetCsrfToken()
			if route.CsrfProtect && ! response.IsValidCsrf() {
				response.SetResponseCode(http.StatusForbidden)
			} else {
				route.Handler(&response)
			}
			fmt.Printf("Handled request: %#v\n", r.URL.Path)
		})
	}
	if len(s.StaticDir) > 0 {
		fmt.Printf("Setting static file handler: %s\n", s.StaticDir)
		s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir(s.StaticDir)))
	}
}

func (s *Server) startServer() {
	srv := &http.Server{
		Handler: s.Router,
		Addr: ":" + strconv.Itoa(s.Port),
	}
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
