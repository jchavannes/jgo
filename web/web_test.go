package web_test

import (
	"github.com/jchavannes/jgo/web"
	"net/http"
)

func ExampleResponse_SetResponseCode() {
	route := web.Route{
		Pattern: "/",
		Handler: func(r *web.Response) {
			r.SetResponseCode(http.StatusNotFound)
		},
	}
	web.Server{
		Routes: []web.Route{
			route,
		},
	}
}
