package health

import (
	"github.com/RichardKnop/go-oauth2-server/util/routes"
	"github.com/gorilla/mux"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
	Close()
}
