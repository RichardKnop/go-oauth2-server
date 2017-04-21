package cmd

import (
	"net/http"
	"time"

	"github.com/adam-hanna/go-oauth2-server/plugins"
	"github.com/adam-hanna/redis-sessions/redis"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
)

// RunServer runs the app
func RunServer(configBackend string) error {
	cnf, db, err := initConfigDB(true, true, configBackend)
	if err != nil {
		return err
	}
	defer db.Close()

	plugins.UseSessionService(redis.NewPluginService())
	if err := plugins.InitServices(cnf, db); err != nil {
		return err
	}
	defer plugins.CloseServices()

	// Start a classic negroni app
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(negroni.NewStatic(http.Dir("public")))

	// Create a router instance
	router := mux.NewRouter()

	// Add routes
	plugins.HealthService.RegisterRoutes(router, "/v1")
	plugins.OauthService.RegisterRoutes(router, "/v1/oauth")
	plugins.WebService.RegisterRoutes(router, "/web")

	// Set the router
	app.UseHandler(router)

	// Run the server on port 8080, gracefully stop on SIGTERM signal
	graceful.Run(":8080", 5*time.Second, app)

	return nil
}
