package routes

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// Route ...
type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []negroni.Handler
}

// AddRoutes adds routes to a router instance. If there are middlewares defined
// for a route, a new negroni app is created and wrapped as a http.Handler
func AddRoutes(routes []Route, router *mux.Router) {
	for _, route := range routes {
		var handler http.Handler
		if len(route.Middlewares) > 0 {
			n := negroni.New()
			for _, middleware := range route.Middlewares {
				n.Use(middleware)
			}
			n.Use(negroni.Wrap(route.HandlerFunc))
			handler = n
		} else {
			handler = route.HandlerFunc
		}

		router.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}
