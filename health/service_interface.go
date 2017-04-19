package health

import (
	"github.com/adam-hanna/go-oauth2-server/util/routes"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	InitHealthService(db *gorm.DB)
	GetRoutes() []routes.Route
	RegisterRoutes(router *mux.Router, prefix string)
}
