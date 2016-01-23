package commands

import (
	"io/ioutil"

	"github.com/AreaHQ/go-fixtures"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/health"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/web"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Migrate migrates the database
func Migrate(db *gorm.DB) error {
	// Bootsrrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := oauth.MigrateAll(db); err != nil {
		return err
	}

	return nil
}

// LoadData loads fixtures
func LoadData(paths []string, cnf *config.Config, db *gorm.DB) error {
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := fixtures.Load(data, db.DB(), cnf.Database.Type); err != nil {
			return err
		}
	}

	return nil
}

// RunServer runs the app
func RunServer(cnf *config.Config, db *gorm.DB) {
	// Initialise the health service
	healthService := health.NewService(db)

	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	// Initialise the web service
	webService := web.NewService(cnf, oauthService)

	// Start a classic negroni app
	app := negroni.Classic()

	// Create a router instance
	router := mux.NewRouter()

	// Add routes for the health service (healthcheck endpoint)
	health.RegisterRoutes(router, healthService)

	// Add routes for the oauth service (REST tokens endpoint)
	oauth.RegisterRoutes(router, oauthService)

	// Add routes for the web package (register, login authorize web pages)
	web.RegisterRoutes(router, webService)

	// Set the router
	app.UseHandler(router)

	// Run the server on port 8080
	app.Run(":8080")
}
