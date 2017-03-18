package web

import (
	"fmt"
	"net/http"
)

func GetFilenameFromRequest(r http.Request) string {
	return r.RequestURI[1:len(r.RequestURI)]
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
