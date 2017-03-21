package web

import (
	"net/http"
	"net/url"
	"time"
	"github.com/jchavannes/jgo/token"
	"fmt"
)

const COOKIE_NAME = "JGoSession"

type Response struct {
	Helper      map[string]interface{}
	Request     Request
	Session     Session
	TemplateDir string
	Writer      http.ResponseWriter
}

func (r *Response) IsValidCsrf() bool {
	requestCsrfToken, err := r.Request.GetCsrfToken()
	return err == nil && requestCsrfToken == r.Session.GetCsrfToken()
}

func (r *Response) ResetOrCreateSession(sessionKey string) {
	r.Session = Session{
		CookieId: token.GetSessionToken(sessionKey),
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

func (r *Response) InitSession(sessionKey string) {
	cookie := r.Request.GetCookie(COOKIE_NAME)
	var validSession bool
	if cookie != "" {
		validSession = token.Validate(cookie, sessionKey)
	}
	if validSession {
		r.Session = Session{
			CookieId: cookie,
		}
	} else {
		r.ResetOrCreateSession(sessionKey)
	}
}

func (r *Response) SetResponseCode(code int) {
	r.Writer.WriteHeader(code)
}

func (r *Response) Write(s string) {
	r.Writer.Write([]byte(s))
}

func (r *Response) Render() {
	requestURI := r.Request.GetURI()
	templateName := requestURI[1:]
	if templateName == "" {
		templateName = "index"
	}
	r.RenderTemplate(templateName)
}

func (r *Response) RenderTemplate(templateName string) {
	renderer, err := GetRenderer(r.TemplateDir)
	if err != nil {
		fmt.Println(err)
	}

	r.Writer.Header().Set("Content-Type", "text/html")

	err = renderer.Render([]string{
		templateName + ".html",
		"404.html",
	}, r.Writer, r.Helper)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Response) SetRedirect(location string) {
	r.Writer.Header().Set("Location", location)
	r.SetResponseCode(http.StatusFound)
}
