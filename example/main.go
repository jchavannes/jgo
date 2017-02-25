package main

import (
	"fmt"
	"github.com/jchavannes/jgo/web"
	"github.com/jchavannes/jgo/example/db"
	"net/http"
)

func main() {
	server := web.Server{
		Port: 8248,
		Sessions: true,
		TemplateDir: "templates",
		StaticDir: "pub",
		Routes: []web.Route{},
	}

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/",
		Handler: func(r *web.Request) {
			r.Render()
		},
	})

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/signup",
		Handler: func(r *web.Request) {
			r.Render()
		},
	})

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/signup-submit",
		CsrfProtect: true,
		Handler: func(r *web.Request) {
			username := r.GetFormValue("username")
			password := r.GetFormValue("password")
			user, err := db.Signup(username, password)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				r.SetResponseCode(http.StatusConflict)
				r.Write("User already exists.")
				return
			}
			fmt.Printf("User: %#v\n", user)
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				return
			}
			session.UserId = user.Id
			err = session.Save()
			if err != nil {
				fmt.Printf("Error saving session: %s\n", err)
				return
			}
		},
	})

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/login",
		Handler: func(r *web.Request) {
			r.Render()
		},
	})

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/login-submit",
		CsrfProtect: true,
		Handler: func(r *web.Request) {
			username := r.GetFormValue("username")
			password := r.GetFormValue("password")
			user, err := db.Login(username, password)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				r.SetResponseCode(http.StatusUnauthorized)
				r.Write(err.Error())
				return
			}
			fmt.Printf("User: %#v\n", user)
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				return
			}
			session.UserId = user.Id
			err = session.Save()
			if err != nil {
				fmt.Printf("Error saving session: %s\n", err)
				return
			}
		},
	})

	server.Routes = append(server.Routes, web.Route{
		Pattern: "/lobby",
		Handler: func(r *web.Request) {
			session, err := db.GetSession(r.Session.CookieId)
			if err != nil {
				r.SetResponseCode(http.StatusUnauthorized)
				r.SetRedirect("/")
				return
			}
			fmt.Printf("Session: %#v\n", session)
			user, err := db.GetUser(db.User{
				Id: session.UserId,
			})
			if err != nil {
				r.SetResponseCode(http.StatusUnauthorized)
				r.SetRedirect("/")
				return
			}
			fmt.Printf("User: %#v\n", user)
			r.Custom = struct{
				Username string
			}{
				Username: user.Username,
			}
			r.Render()
		},
	})

	server.Run()
}
