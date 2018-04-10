package web

type Route struct {
	Pattern     string
	Name        string // Optional
	Handler     func(*Response)
	CsrfProtect bool
	NeedsLogin  bool
}

func Routes(routeSets ...[]Route) []Route {
	var allRoutes []Route
	for _, routeSet := range routeSets {
		allRoutes = append(allRoutes, routeSet...)
	}
	return allRoutes
}
