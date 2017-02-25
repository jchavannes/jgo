package web

import (
	"net/http"
	"net/url"
	"time"
)

const COOKIE_NAME = "JGoSession"

type Request struct {
	HttpResponseWriter http.ResponseWriter
	HttpRequest        http.Request
	TemplateDir        string
	Session            Session
	Custom             interface{}
}

func (r *Request) IsCsrfPresentAndValid() bool {
	if r.HttpRequest.Method != "POST" {
		return false
	}
	csrfToken := r.HttpRequest.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		return false
	}
	return csrfToken == r.Session.GetCsrfToken()
}

func (r *Request) GetFormValue(key string) string {
	r.HttpRequest.ParseForm()
	return r.HttpRequest.Form.Get(key)
}

func (r *Request) GetHeader(key string) string {
	return r.HttpRequest.Header.Get(key)
}

func (r *Request) GetBaseUrl() string {
	apppath := r.GetHeader("Apppath")
	if apppath == "" {
		apppath = "/"
	}
	return apppath
}

func (r *Request) ResetOrCreateSession() {
	r.Session = Session{
		CookieId: CreateToken(),
	}
	cookie := http.Cookie{
		Name: COOKIE_NAME,
		Value: url.QueryEscape(r.Session.CookieId),
		Path: "/",
		HttpOnly: true,
		MaxAge: int(time.Hour) * 24 * 30,
	}
	http.SetCookie(r.HttpResponseWriter, &cookie)
}

func (r *Request) InitSession() {
	cookie, err := r.HttpRequest.Cookie(COOKIE_NAME)
	if err != nil || cookie.Value == "" {
		r.ResetOrCreateSession()
	} else {
		r.Session = Session{
			CookieId: cookie.Value,
		}
	}
}

func (r *Request) SetResponseCode(code int) {
	r.HttpResponseWriter.WriteHeader(code)
}

func (r *Request) Write(s string) {
	r.HttpResponseWriter.Write([]byte(s))
}

func (r *Request) Render() {
	requestURI := r.HttpRequest.RequestURI
	templateName := requestURI[1:]
	if templateName == "" {
		templateName = "index"
	}
	r.RenderTemplate(templateName)
}

func (r *Request) RenderTemplate(templateName string) {
	renderer, err := GetRenderer(r.TemplateDir)
	check(err)

	r.HttpResponseWriter.Header().Set("Content-Type", "text/html")

	renderer.Render([]string{
		templateName + ".html",
		"404.html",
	}, r.HttpResponseWriter, r)
}

func (r *Request) SetRedirect(location string) {
	r.HttpResponseWriter.Header().Set("Location", location)
}
