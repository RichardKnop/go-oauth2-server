package cmd

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/RichardKnop/go-oauth2-server/health"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/web"
	"github.com/jinzhu/gorm"
)

var (
	healthService health.ServiceInterface
	oauthService  oauth.ServiceInterface
	webService    web.ServiceInterface
)

// initConfigDB loads the configuration and connects to the database
func initConfigDB(mustLoadOnce, keepReloading bool) (*config.Config, *gorm.DB, error) {
	// Config
	cnf := config.NewConfig(mustLoadOnce, keepReloading)

	// Database
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return nil, nil, err
	}

	return cnf, db, nil
}

// initServices starts up all services and sets above defined variables
func initServices(cnf *config.Config, db *gorm.DB) error {
	healthService = health.NewService(db)

	oauthService = oauth.NewService(cnf, db)

	webService = web.NewService(cnf, oauthService)

	return nil
}
