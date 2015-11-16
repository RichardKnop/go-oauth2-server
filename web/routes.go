package web

import (
	"github.com/RichardKnop/go-oauth2-server/routes"

	"github.com/codegangsta/negroni"
)

// Routes for web pages
var Routes = []routes.Route{
	routes.Route{
		Name:        "register_form",
		Methods:     []string{"GET"},
		Pattern:     "/register",
		HandlerFunc: registerForm,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(guestMiddleware),
		},
	},
	routes.Route{
		Name:        "register",
		Methods:     []string{"POST"},
		Pattern:     "/register",
		HandlerFunc: register,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(guestMiddleware),
		},
	},
	routes.Route{
		Name:        "login_form",
		Methods:     []string{"GET"},
		Pattern:     "/login",
		HandlerFunc: loginForm,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(guestMiddleware),
		},
	},
	routes.Route{
		Name:        "login",
		Methods:     []string{"POST"},
		Pattern:     "/login",
		HandlerFunc: login,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(guestMiddleware),
		},
	},
	routes.Route{
		Name:        "logout",
		Methods:     []string{"GET"},
		Pattern:     "/logout",
		HandlerFunc: logout,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(loggedInMiddleware),
		},
	},
	routes.Route{
		Name:        "authorize_form",
		Methods:     []string{"GET"},
		Pattern:     "/authorize",
		HandlerFunc: authorizeForm,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(loggedInMiddleware),
			new(clientMiddleware),
		},
	},
	routes.Route{
		Name:        "authorize",
		Methods:     []string{"POST"},
		Pattern:     "/authorize",
		HandlerFunc: authorize,
		Middlewares: []negroni.Handler{
			new(parseFormMiddleware),
			new(loggedInMiddleware),
			new(clientMiddleware),
		},
	},
}
