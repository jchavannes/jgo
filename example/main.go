package main

import (
	"github.com/jchavannes/jgo/web"
)

func main() {
	server := web.Server{
		Port: 8080,
		Template: web.Template{
			Pattern: "/",
			Directory: "./",
		},
		Static: web.Static{
			Pattern: "/pub",
			Directory: "pub",
		},
		Routes: []web.Route{{
			Pattern: "/post",
			CsrfProtect: true,
			Handler: func(r *web.Request) {
				r.Write("CSRF protected")
			},
		}},
	}
	server.Run()
}
