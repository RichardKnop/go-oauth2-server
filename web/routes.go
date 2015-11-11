package web

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// Routes for web pages
var Routes = []routes.Route{
	routes.Route{
		Name:        "Register",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/register",
		HandlerFunc: register,
	},
	routes.Route{
		Name:        "Login",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/login",
		HandlerFunc: login,
	},
	routes.Route{
		Name:        "Authorize",
		Methods:     []string{"GET"},
		Pattern:     "/authorize",
		HandlerFunc: authorize,
	},
}
