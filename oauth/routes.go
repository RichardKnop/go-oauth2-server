package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// Routes for the oauth service
var Routes = []routes.Route{
	routes.Route{
		Name:        "oauth_tokens",
		Methods:     []string{"POST"},
		Pattern:     "/tokens",
		HandlerFunc: handleTokens,
	},
}
