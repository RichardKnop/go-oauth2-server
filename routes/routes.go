package routes

import (
	"net/http"

	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

// Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []negroni.Handler
}

// AddRoutes adds routes to a router instance. If there are middlewares defined
// for a route, a new negroni app is created and wrapped as a http.Handler
func AddRoutes(routes []Route, router *mux.Router) {
	var (
		handler http.Handler
		n       *negroni.Negroni
	)

	for _, route := range routes {
		// Add any specified middlewares
		if len(route.Middlewares) > 0 {
			n = negroni.New()

			for _, middleware := range route.Middlewares {
				n.Use(middleware)
			}

			// Wrap the handler in the negroni app with middlewares
			n.Use(negroni.Wrap(route.HandlerFunc))
			handler = n
		} else {
			handler = route.HandlerFunc
		}

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}
