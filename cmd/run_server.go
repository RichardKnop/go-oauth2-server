package cmd

import (
	"net/http"
	"time"

	"github.com/adam-hanna/go-oauth2-server/services"
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

	// configure redis for session store
	sessionSecrets := make([][]byte, 1)
	sessionSecrets[0] = []byte(cnf.Session.Secret)
	redisConfig := redis.ConfigType{
		Size:           10,
		Network:        "tcp",
		Address:        ":6379",
		Password:       "",
		SessionSecrets: sessionSecrets,
	}

	// start the services
	services.UseSessionService(redis.NewService(cnf, redisConfig))
	if err := services.InitServices(cnf, db); err != nil {
		return err
	}
	defer services.CloseServices()

	// Start a classic negroni app
	app := negroni.New()
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(gzip.Gzip(gzip.DefaultCompression))
	app.Use(negroni.NewStatic(http.Dir("public")))

	// Create a router instance
	router := mux.NewRouter()

	// Add routes
	services.HealthService.RegisterRoutes(router, "/v1")
	services.OauthService.RegisterRoutes(router, "/v1/oauth")
	services.WebService.RegisterRoutes(router, "/web")

	// Set the router
	app.UseHandler(router)

	// Run the server on port 8080, gracefully stop on SIGTERM signal
	graceful.Run(":8080", 5*time.Second, app)

	return nil
}
