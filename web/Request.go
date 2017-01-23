package web

import "net/http"

type Request struct {
	HttpResponseWriter *http.ResponseWriter
	HttpRequest        *http.Request
}

func (r *Request) Write(s string) {
	writer := *r.HttpResponseWriter
	writer.Write([]byte(s))
}
