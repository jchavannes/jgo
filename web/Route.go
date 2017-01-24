package web

type Route struct {
	Pattern     string
	Handler     func(*Request)
	CsrfProtect bool
}
