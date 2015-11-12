package web

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// Routes for web pages
var Routes = []routes.Route{
	routes.Route{
		Name:        "register_form",
		Methods:     []string{"GET"},
		Pattern:     "/register",
		HandlerFunc: registerForm,
	},
	routes.Route{
		Name:        "register",
		Methods:     []string{"POST"},
		Pattern:     "/register",
		HandlerFunc: register,
	},
	routes.Route{
		Name:        "login_form",
		Methods:     []string{"GET"},
		Pattern:     "/login",
		HandlerFunc: loginForm,
	},
	routes.Route{
		Name:        "login",
		Methods:     []string{"POST"},
		Pattern:     "/login",
		HandlerFunc: login,
	},
	routes.Route{
		Name:        "authorize_form",
		Methods:     []string{"GET"},
		Pattern:     "/authorize",
		HandlerFunc: authorizeForm,
	},
	routes.Route{
		Name:        "authorize",
		Methods:     []string{"POST"},
		Pattern:     "/authorize",
		HandlerFunc: authorize,
	},
}
