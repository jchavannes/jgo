package web

import (
	"net/http"
)

type Handler struct {
	Handler func(http.ResponseWriter, *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler(w, r)
}
