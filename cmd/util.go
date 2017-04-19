package cmd

import (
	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/database"
	"github.com/adam-hanna/go-oauth2-server/health"
	"github.com/adam-hanna/go-oauth2-server/oauth"
	"github.com/adam-hanna/go-oauth2-server/plugins"
	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/adam-hanna/go-oauth2-server/web"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

var (
	healthService  health.ServiceInterface
	oauthService   oauth.ServiceInterface
	webService     web.ServiceInterface
	sessionService session.ServiceInterface
)

// initConfigDB loads the configuration and connects to the database
func initConfigDB(mustLoadOnce, keepReloading bool, configBackend string) (*config.Config, *gorm.DB, error) {
	// Config
	cnf := config.NewConfig(mustLoadOnce, keepReloading, configBackend)

	// Database
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return nil, nil, err
	}

	return cnf, db, nil
}

// initServices starts up all services and sets above defined variables
func initServices(cnf *config.Config, db *gorm.DB) error {
	healthService = plugins.NewHealthService(db)
	if healthService == (*plugins.CustomHealthService)(nil) {
		healthService = health.NewService(db)
	}

	oauthService = plugins.NewOauthService(cnf, db)
	if oauthService == (*plugins.CustomAuthService)(nil) {
		oauthService = oauth.NewService(cnf, db)
	}

	sessionService = plugins.NewSessionService(cnf)
	if sessionService == (*plugins.CustomSessionService)(nil) {
		// note: default session store is CookieStore
		sessionService = session.NewService(cnf, sessions.NewCookieStore([]byte(cnf.Session.Secret)))
	}

	webService = plugins.NewWebService(cnf, oauthService, sessionService)
	if webService == (*plugins.CustomWebService)(nil) {
		webService = web.NewService(cnf, oauthService, sessionService)
	}

	return nil
}
