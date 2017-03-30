package web

type Route struct {
	Pattern     string
	Name        string // Optional
	Handler     func(*Response)
	CsrfProtect bool
}
