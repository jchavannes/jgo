package main

import (
	"fmt"
	"github.com/jchavannes/jgo/example/auth"
	"github.com/jchavannes/jgo/example/db"
	"github.com/jchavannes/jgo/web"
	"net/http"
)

var (
	defaultRoute = web.Route{
		Pattern: "/",
		Handler: func(r *web.Request) {
			r.Render()
		},
	}

	signupRoute = web.Route{
		Pattern: "/signup",
		Handler: func(r *web.Request) {
			r.Render()
		},
	}

	signupSubmitRoute = web.Route{
		Pattern: "/signup-submit",
		CsrfProtect: true,
		Handler: func(r *web.Request) {
			username := r.GetFormValue("username")
			password := r.GetFormValue("password")
			user, err := auth.Signup(username, password)
			if err != nil {
				fmt.Printf("Error signing up: %s\n", err)
				r.SetResponseCode(http.StatusConflict)
				r.Write("User already exists.")
				return
			}
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}
			session.UserId = user.Id
			err = session.Save()
			if err != nil {
				fmt.Printf("Error saving session: %s\n", err)
				return
			}
		},
	}

	loginRoute = web.Route{
		Pattern: "/login",
		Handler: func(r *web.Request) {
			r.Render()
		},
	}

	loginSubmitRoute = web.Route{
		Pattern: "/login-submit",
		CsrfProtect: true,
		Handler: func(r *web.Request) {
			username := r.GetFormValue("username")
			password := r.GetFormValue("password")
			user, err := auth.Login(username, password)
			if err != nil {
				fmt.Printf("Error logging in: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.Write(err.Error())
				return
			}
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				return
			}
			session.UserId = user.Id
			err = session.Save()
			if err != nil {
				fmt.Printf("Error saving session: %s\n", err)
				return
			}
		},
	}

	lobbyRoute = web.Route{
		Pattern: "/lobby",
		Handler: func(r *web.Request) {
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				fmt.Printf("Error getting session: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.SetRedirect("/")
				return
			}
			user, err := db.GetUser(db.User{
				Id: session.UserId,
			})
			if err != nil {
				fmt.Printf("Error getting user: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.SetRedirect("/")
				return
			}
			r.Helper["Username"] = user.Username
			r.Render()
		},
	}

	initRequest = func(r *web.Request) {
		baseUrl := r.GetHeader("AppPath")
		if baseUrl == "" {
			baseUrl = "/"
		}
		r.Helper["BaseUrl"] = baseUrl
	}
)

func main() {
	server := web.Server{
		Port: 8248,
		Sessions: true,
		TemplateDir: "templates",
		StaticDir: "pub",
		InitRequest: initRequest,
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
