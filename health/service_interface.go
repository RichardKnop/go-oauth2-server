package health

import (
	"github.com/gorilla/mux"
	"github.com/RichardKnop/go-oauth2-server/routes"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
}
