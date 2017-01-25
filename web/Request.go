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
	Session            Session
}

func (r *Request) IsCsrfPresentAndValid() bool {
	if r.HttpRequest.Method != "POST" {
		return false
	}
	csrfToken := r.GetFormValue("csrf-token")
	if csrfToken == "" {
		return false
	}
	return csrfToken == r.Session.GetCsrfToken()
}

func (r *Request) GetFormValue(key string) string {
	for k, v := range r.HttpRequest.Form {
		if k == key {
			return v[0]
		}
	}
	return ""
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

func (r *Request) Write(s string) {
	r.HttpResponseWriter.Write([]byte(s))
}
