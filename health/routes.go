package health

import (
	"github.com/RichardKnop/go-oauth2-server/routes"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers route handlers for the health service
func RegisterRoutes(router *mux.Router, service *Service) {
	routes.AddRoutes(newRoutes(service), router)
}

// newRoutes returns []routes.Route slice for the health service
func newRoutes(service *Service) []routes.Route {
	return []routes.Route{
		routes.Route{
			Name:        "health_check",
			Method:      "GET",
			Pattern:     "/v1/health",
			HandlerFunc: service.healthcheck,
		},
	}
}
