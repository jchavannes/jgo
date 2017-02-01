package main

import (
	"github.com/jchavannes/jgo/web"
)

func main() {
	server := web.Server{
		Port: 8080,
		Sessions: true,
		TemplateDir: "templates",
		StaticDir: "pub",
		Routes: []web.Route{{
			Pattern: "/",
			Handler: func(r *web.Request) {
				r.Render("index")
			},
		}, {
			Pattern: "/signup",
			Handler: func(r *web.Request) {
				r.Render("signup")
			},
		}, {
			Pattern: "/post",
			CsrfProtect: true,
			Handler: func(r *web.Request) {
				r.Write("CSRF protected")
			},
		}},
	}
	server.Run()
}
