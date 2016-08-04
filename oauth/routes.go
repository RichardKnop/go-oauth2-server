package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the oauth service
func RegisterRoutes(router *mux.Router, service ServiceInterface) {
	subRouter := router.PathPrefix("/v1/oauth").Subrouter()
	routes.AddRoutes(newRoutes(service), subRouter)
}

// newRoutes returns []routes.Route slice for the oauth service
func newRoutes(service ServiceInterface) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "oauth_tokens",
			Method:      "POST",
			Pattern:     "/tokens",
			HandlerFunc: service.tokensHandler,
		},
		routes.Route{
			Name:        "oauth_introspect",
			Method:      "POST",
			Pattern:     "/introspect",
			HandlerFunc: service.introspectHandler,
		},
	}
}
