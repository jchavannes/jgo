package main

import (
	"github.com/jchavannes/jgo/web"
)

func main() {
	server := web.Server{ Port: 8080, TemplateDirectory: "./", }
	server.Run()
}
