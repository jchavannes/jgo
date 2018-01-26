package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/token"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response objects are passed to handlers to respond to requests.
// Includes abstracted access to request and session information
type Response struct {
	Helper  map[string]interface{}
	Request Request
	Server  *Server
	Session Session
	rcSet   bool
	Writer  http.ResponseWriter
}

// Checks that CSRF token in request matches one for session.
// Tokens are kept in memory and do not persist between instances or restarts.
func (r *Response) IsValidCsrf() bool {
	requestCsrfToken, err := r.Request.GetCsrfToken()
	return err == nil && requestCsrfToken == r.Session.GetCsrfToken()
}

// Sets a new session cookie.
func (r *Response) ResetOrCreateSession() {
	r.Session = Session{
		CookieId: token.GetSessionToken(r.Server.SessionKey),
	}
	cookie := http.Cookie{
		Name:     r.Server.GetCookieName(),
		Value:    url.QueryEscape(r.Session.CookieId),
		Path:     "/",
		HttpOnly: true,
		Secure:   !r.Server.InsecureCookie,
		MaxAge:   int(time.Hour) * 24 * 30,
	}
	http.SetCookie(r.Writer, &cookie)
}

// Either gets existing session token or creates a new one.
func (r *Response) InitSession() {
	cookie := r.Request.GetCookie(r.Server.GetCookieName())
	var validSession bool
	if cookie != "" {
		validSession = token.Validate(cookie, r.Server.SessionKey)
	}
	if validSession {
		r.Session = Session{
			CookieId: cookie,
		}
	} else {
		r.ResetOrCreateSession()
	}
}

func (r *Response) SetResponseCode(code int) error {
	if r.rcSet {
		return errors.New("response code already set")
	}
	r.Writer.WriteHeader(code)
	r.rcSet = true
	return nil
}

func (r *Response) SetHeader(key string, value string) {
	r.Writer.Header().Set(key, value)
}

func (r *Response) Write(s string) {
	r.Writer.Write([]byte(s))
}

func (r *Response) WriteJson(i interface{}, pretty bool) {
	var text []byte
	if pretty {
		text, _ = json.MarshalIndent(i, "", "  ")
	} else {
		text, _ = json.Marshal(i)
	}
	r.SetHeader("Content-Type", "application/json")
	r.Write(string(text))
}

func (r *Response) Render() error {
	requestURI := r.Request.GetURI()
	templateName := strings.Split(requestURI[1:], "?")[0]
	return r.RenderTemplate(templateName)
}

func (r *Response) RenderTemplate(templateName string) error {
	renderer, err := GetRenderer(r.Server.TemplatesDir)
	if err != nil {
		fmt.Println(err)
	}

	r.SetHeader("Content-Type", "text/html")

	if strings.HasPrefix(templateName, "/") {
		templateName = templateName[1:]
	}

	var indexTemplate string
	if templateName == "" {
		indexTemplate = "index.html"
	} else {
		indexTemplate = templateName + "/index.html"
	}

	err = renderer.Render([]string{
		templateName + ".html",
		indexTemplate,
		"404.html",
	}, r.Writer, r.Helper)

	if err != nil {
		fmt.Printf("Error rendering template: %s\n", err)
	}

	return err
}

func (r *Response) SetRedirect(location string) {
	r.SetHeader("Location", location)
	r.SetResponseCode(http.StatusFound)
}

func (r *Response) GetWebSocket() (*Socket, error) {
	upgrader := websocket.Upgrader{
		EnableCompression: true,
	}
	conn, err := upgrader.Upgrade(r.Writer, &r.Request.HttpRequest, nil)
	if err != nil {
		return nil, err
	}
	return &Socket{ws: conn}, nil
}

func (r *Response) Error(err error, responseCode int) {
	r.SetResponseCode(responseCode)
	fmt.Println(jerr.Get(fmt.Sprintf("Error with request: %#v", r.Request.HttpRequest.URL.Path), err))
}

func (r *Response) LogComplete() {
	fmt.Printf("Handled request: %#v\n", r.Request.HttpRequest.URL.Path)
}
