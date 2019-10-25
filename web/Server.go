package web

import (
	"github.com/gorilla/mux"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Name of cookie used for sessions.
const CookieName = "session_id"

// Maximum number of routes to display when starting up.
const MaxRoutesDisplay = 3

type Server struct {
	AllowedExtensions []string
	IsLoggedIn        func(*Response) bool
	NotFoundHandler   func(*Response)
	Port              int
	PreHandler        func(*Response)
	PostHandler       func(*Response)
	ErrorHandler      func(*Response, error)
	GetCsrfToken      func(string) string
	Routes            []Route
	SessionKey        string
	StaticFilesDir    string
	StrictSlash       bool
	TemplatesDir      string
	UseAutoRender     bool
	UseSessions       bool
	InsecureCookie    bool
	CookiePrefix      string
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
	s.router.SkipClean(true)
	return s.startServer()
}

func (s *Server) addCatchAllRoute() {
	if len(s.StaticFilesDir) > 0 {
		jlog.Logf("Static files directory: %s\n", s.StaticFilesDir)
	}
	if len(s.TemplatesDir) > 0 {
		jlog.Logf("Templates directory: %s\n", s.TemplatesDir)
	}
	s.router.PathPrefix("/").Handler(Handler{
		Handler: func(w http.ResponseWriter, r *http.Request) {
			response := getResponse(w, r, s, true)
			defer response.LogComplete()
			if response.ResponseCodeSet() {
				return
			}
			if len(s.StaticFilesDir) > 0 {
				allowedFileTypes := GetDefaultAllowedFileExtensions()
				if len(s.AllowedExtensions) > 0 {
					allowedFileTypes = s.AllowedExtensions
				}
				for _, fileType := range allowedFileTypes {
					if strings.HasSuffix(response.Request.HttpRequest.URL.Path, "."+fileType) || fileType == "*" {
						if fileType == "css" {
							var re = regexp.MustCompile(`(.*)-[\d]+(\.css)`)
							path := re.ReplaceAllString(response.Request.HttpRequest.URL.Path, `$1$2`)
							http.ServeFile(w, r, s.StaticFilesDir+path)
							return
						} else if fileType == "js" {
							var re = regexp.MustCompile(`(.*)-[\d]+(\.js)`)
							path := re.ReplaceAllString(response.Request.HttpRequest.URL.Path, `$1$2`)
							http.ServeFile(w, r, s.StaticFilesDir+path)
							return
						}
						http.FileServer(http.Dir(s.StaticFilesDir)).ServeHTTP(response.Writer, &response.Request.HttpRequest)
						return
					}
				}
			}

			if len(s.TemplatesDir) > 0 && s.UseAutoRender {
				templateName := response.Request.GetPotentialFilename()

				if templateName == "" {
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
			if s.PostHandler != nil {
				s.PostHandler(&response)
			}
		},
	})
}

func (s *Server) setupHandlers() {
	s.router = mux.NewRouter()
	s.router.StrictSlash(s.StrictSlash)
	var showEachRoute = len(s.Routes) <= MaxRoutesDisplay
	if ! showEachRoute {
		jlog.Logf("Adding %d patterns to router.\n", len(s.Routes))
	}
	for _, routeTemp := range s.Routes {
		route := routeTemp
		name := ""
		if len(route.Name) > 0 {
			name = " (" + route.Name + ")"
		}
		if showEachRoute {
			jlog.Logf("Adding pattern to router: %s%s\n", route.Pattern, name)
		}
		s.router.HandleFunc(route.Pattern, func(w http.ResponseWriter, r *http.Request) {
			response := getResponse(w, r, s, false)
			response.Pattern = route.Pattern
			defer response.LogComplete()
			if response.ResponseCodeSet() {
				return
			}
			if route.CsrfProtect && ! response.IsValidCsrf() {
				response.SetResponseCode(http.StatusForbidden)
				return
			}
			if s.IsLoggedIn != nil && route.NeedsLogin && ! s.IsLoggedIn(&response) {
				return
			}
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						response.Error(jerr.Getf(err, "fatal error processing response, stack: %s", debug.Stack()), http.StatusInternalServerError)
					}
				}
			}()
			route.Handler(&response)
			if s.PostHandler != nil {
				s.PostHandler(&response)
			}
		})
	}
}

func (s *Server) GetCookieName() string {
	cookieName := CookieName
	if s.CookiePrefix != "" {
		cookieName = s.CookiePrefix + "_" + cookieName
	}
	return cookieName
}

func getResponse(w http.ResponseWriter, r *http.Request, s *Server, static bool) Response {
	response := Response{
		Helper:  make(map[string]interface{}),
		Request: Request{HttpRequest: *r},
		Server:  s,
		StartTs: time.Now(),
		Writer:  w,
		Static:  static,
	}
	response.Helper["URI"] = r.RequestURI
	if s.UseSessions {
		response.InitSession()
		response.Helper["CsrfToken"] = response.GetCsrfToken()
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
	jlog.Logf("Starting server on port %d...\n", s.Port)
	return srv.ListenAndServe()
}
