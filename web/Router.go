package web

type Router struct {
	NotFoundHandler func(*Response)
	PatternPrefix   string
	PreHandler      func(*Response)
	Routes          []Route
}

func (r Router) GetPatternPrefix() string {
	if r.PatternPrefix == "" {
		return "/"
	}
	return r.PatternPrefix
}

func (r Router) GetRoutePattern(route Route) string {
	return r.PatternPrefix + route.Pattern
}
