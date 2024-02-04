package web_test

import (
	"github.com/jchavannes/jgo/web"
	"net/http"
)

func ExampleResponse_SetResponseCode() {
	_ = web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			r.SetResponseCode(http.StatusNotFound)
		},
	}
}

func ExampleCreateToken() {
	token := web.CreateToken()
	println(token)
}
