package oauth2

import (
	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

// NewRoutes - returns routes for the main app
func NewRoutes(cnf *config.Config, db *gorm.DB) []*rest.Route {
	return []*rest.Route{
		rest.Post("/api/v1/users", func(w rest.ResponseWriter, r *rest.Request) {
			register(w, r, cnf, db)
		}),
		rest.Post("/api/v1/tokens", func(w rest.ResponseWriter, r *rest.Request) {
			tokensHandler(w, r, cnf, db)
		}),
	}
}
