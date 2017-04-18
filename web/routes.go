package web

import (
	"github.com/adam-hanna/go-oauth2-server/util/routes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// RegisterRoutes registers route handlers for the health service
func (s *Service) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	routes.AddRoutes(s.GetRoutes(), subRouter)
}

// GetRoutes returns []routes.Route slice for the health service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "register_form",
			Method:      "GET",
			Pattern:     "/register",
			HandlerFunc: s.registerForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(s),
				newClientMiddleware(s),
			},
		},
		routes.Route{
			Name:        "register",
			Method:      "POST",
			Pattern:     "/register",
			HandlerFunc: s.register,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(s),
				newClientMiddleware(s),
			},
		},
		routes.Route{
			Name:        "login_form",
			Method:      "GET",
			Pattern:     "/login",
			HandlerFunc: s.loginForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(s),
				newClientMiddleware(s),
			},
		},
		routes.Route{
			Name:        "login",
			Method:      "POST",
			Pattern:     "/login",
			HandlerFunc: s.login,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newGuestMiddleware(s),
				newClientMiddleware(s),
			},
		},
		routes.Route{
			Name:        "logout",
			Method:      "GET",
			Pattern:     "/logout",
			HandlerFunc: s.logout,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(s),
			},
		},
		routes.Route{
			Name:        "authorize_form",
			Method:      "GET",
			Pattern:     "/authorize",
			HandlerFunc: s.authorizeForm,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(s),
				newClientMiddleware(s),
			},
		},
		routes.Route{
			Name:        "authorize",
			Method:      "POST",
			Pattern:     "/authorize",
			HandlerFunc: s.authorize,
			Middlewares: []negroni.Handler{
				new(parseFormMiddleware),
				newLoggedInMiddleware(s),
				newClientMiddleware(s),
			},
		},
	}
}
