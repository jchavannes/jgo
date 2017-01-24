package web

import (
	"net/http"
	"log"
)

func GetFilenameFromRequest(r http.Request) string {
	return r.RequestURI[1:len(r.RequestURI)]
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
