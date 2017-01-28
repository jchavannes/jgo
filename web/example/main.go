package main

import (
	"github.com/jchavannes/jgo/web"
)

func main() {
	server := web.Server{
		Port: 8080,
		TemplateDirectory: "./",
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
