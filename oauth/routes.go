package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

// NewRoutes returns routes for OAuth service
func NewRoutes(cnf *config.Config, db *gorm.DB) []*rest.Route {
	s := &service{cnf: cnf, db: db}
	return []*rest.Route{
		rest.Post("/oauth2/api/v1/tokens", s.handleTokens),
	}
}
