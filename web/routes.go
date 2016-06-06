package web

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
	"github.com/urfave/negroni"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the web package
func RegisterRoutes(router *mux.Router, service *Service) {
	subRouter := router.PathPrefix("/web").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the web package
func newRoutes(service *Service) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "register_form",
			Method:      "GET",
			Pattern:     "/register",
			HandlerFunc: service.registerForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "register",
			Method:      "POST",
			Pattern:     "/register",
			HandlerFunc: service.register,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "login_form",
			Method:      "GET",
			Pattern:     "/login",
			HandlerFunc: service.loginForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: service.login,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "logout",
			Method:      "GET",
			Pattern:     "/logout",
			HandlerFunc: service.logout,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
			},
		},
		routes.Route{
			Name:        "authorize_form",
			Method:      "GET",
			Pattern:     "/authorize",
			HandlerFunc: service.authorizeForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
				newClientMiddleware(service),
			},
		},
		routes.Route{
			Name:        "authorize",
			Method:      "POST",
			Pattern:     "/authorize",
			HandlerFunc: service.authorize,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(service),
				newClientMiddleware(service),
			},
		},
	}
}
