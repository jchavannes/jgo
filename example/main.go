package main

import (
	"fmt"
	"github.com/jchavannes/jgo/example/db"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var (
	defaultRoute = web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			session, err := db.GetSession(r.Session.CookieId)
			if session.IsLoggedIn() {
				r.SetRedirect(getBaseUrl(r.Request) + "lobby")
			} else {
				r.Render()
			}
		},
	}

	signupRoute = web.Route{
		Pattern: "/signup",
		Handler: func(r *web.Response) {
			r.Render()
		},
	}

	signupSubmitRoute = web.Route{
		Pattern: "/signup-submit",
		CsrfProtect: true,
		Handler: func(r *web.Response) {
			username := r.Request.GetFormValue("username")
			password := r.Request.GetFormValue("password")

			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}

			_, err = session.Signup(username, password)
			if err != nil {
				fmt.Printf("Error signing up: %s\n", err)
				r.SetResponseCode(http.StatusConflict)
				return
			}
		},
	}

	loginRoute = web.Route{
		Pattern: "/login",
		Handler: func(r *web.Response) {
			r.Render()
		},
	}

	loginSubmitRoute = web.Route{
		Pattern: "/login-submit",
		CsrfProtect: true,
		Handler: func(r *web.Response) {
			username := r.Request.GetFormValue("username")
			password := r.Request.GetFormValue("password")

			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}

			_, err = session.Login(username, password)
			if err != nil {
				fmt.Printf("Error logging in: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.Write(err.Error())
			}
		},
	}

	lobbyRoute = web.Route{
		Pattern: "/lobby",
		Handler: func(r *web.Response) {
			session, err := db.GetSession(r.Session.CookieId)
			if err == nil {
				user, err := db.GetUser(db.User{
					Id: session.UserId,
				})
				if err == nil {
					r.Helper["Username"] = user.Username
					r.Render()
					return
				}
			}
			fmt.Printf("Error getting user: %s\n", err)
			r.SetResponseCode(http.StatusUnauthorized)
			r.SetRedirect("/")
		},
	}

	initRequest = func(r *web.Response) {
		r.Helper["BaseUrl"] = getBaseUrl(r.Request)
	}

	getBaseUrl = func(r web.Request) string {
		baseUrl := r.GetHeader("AppPath")
		if baseUrl == "" {
			baseUrl = "/"
		}
		return baseUrl
	}
)

func main() {
	server := web.Server{
		Port: 8248,
		Sessions: true,
		TemplateDir: "templates",
		StaticDir: "pub",
		InitResponse: initRequest,
		Routes: []web.Route{
			defaultRoute,
			signupRoute,
			signupSubmitRoute,
			loginRoute,
			loginSubmitRoute,
			lobbyRoute,
		},
	}

	server.Run()
}
