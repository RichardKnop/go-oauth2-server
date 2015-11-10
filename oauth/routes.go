package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// Routes for the oauth service
var Routes = []routes.Route{
	routes.Route{
		Name:        "Tokens",
		Methods:     []string{"POST"},
		Pattern:     "/api/v1/tokens",
		HandlerFunc: handleTokens,
	},
}
