package web

import (
	"net/http"
	"net/url"
	"time"
	"github.com/jchavannes/jgo/token"
	"fmt"
	"github.com/gorilla/websocket"
	"encoding/json"
	"errors"
)

// Name of cookie used for sessions.
const COOKIE_NAME = "JGoSession"

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
		Name: COOKIE_NAME,
		Value: url.QueryEscape(r.Session.CookieId),
		Path: "/",
		HttpOnly: true,
		MaxAge: int(time.Hour) * 24 * 30,
	}
	http.SetCookie(r.Writer, &cookie)
}

// Either gets existing session token or creates a new one.
func (r *Response) InitSession() {
	cookie := r.Request.GetCookie(COOKIE_NAME)
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
		return errors.New("Response code already set")
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
	r.Write(string(text))
}

func (r *Response) Render() error {
	requestURI := r.Request.GetURI()
	templateName := requestURI[1:]
	if templateName == "" {
		templateName = "index"
	}
	return r.RenderTemplate(templateName)
}

func (r *Response) RenderTemplate(templateName string) error {
	renderer, err := GetRenderer(r.Server.TemplatesDir)

	if err != nil {
		fmt.Println(err)
	}

	r.SetHeader("Content-Type", "text/html")

	err = renderer.Render([]string{
		templateName + ".html",
		"404.html",
	}, r.Writer, r.Helper)

	return err
}

func (r *Response) SetRedirect(location string) {
	r.SetHeader("Location", location)
	r.SetResponseCode(http.StatusFound)
}

func (r *Response) GetWebSocket() (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(r.Writer, &r.Request.HttpRequest, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (r *Response) LogComplete() {
	fmt.Printf("Handled request: %#v\n", r.Request.HttpRequest.URL.Path)
}
