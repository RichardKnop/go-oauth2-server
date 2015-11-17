package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RichardKnop/go-oauth2-server/accounts"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/web"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	cliApp *cli.App
	webApp *negroni.Negroni
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "go-oauth2-server"
	cliApp.Usage = "OAuth 2.0 Server"
	cliApp.Author = "Richard Knop"
	cliApp.Email = "risoknop@gmail.com"
	cliApp.Version = "0.0.0"
}

func main() {
	// Load the configuration, connect to the database
	cnf := config.NewConfig()
	db, err := database.NewDatabase(cnf)
	if err != nil {
		log.Fatal(err)
	}
	// Disable logging
	db.LogMode(false)

	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:   "migrate",
			Usage:  "run migrations",
			Action: func(c *cli.Context) { migrate(db) },
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) {
				runServer(cnf, db)
			},
		},
	}

	cliApp.Run(os.Args)
}

func migrate(db *gorm.DB) {
	// Bootsrrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		log.Fatal(err)
	}
	// Run migrations for the oauth service
	if err := oauth.MigrateAll(db); err != nil {
		log.Fatal(err)
	}
}

func runServer(cnf *config.Config, db *gorm.DB) {
	// Initialise the oauth service
	oauthService := oauth.NewService(cnf, db)

	// Initialise the accounts service
	_ = accounts.NewService(cnf, db, oauthService)

	// Initialise the web service
	_ = web.NewService(cnf, oauthService)

	// Start a classic negroni app
	webApp := negroni.Classic()

	// Create a router instance
	router := mux.NewRouter().StrictSlash(true)

	var subRouter *mux.Router

	// Add routes for the oauth service
	subRouter = router.PathPrefix("/oauth/api/v1").Subrouter()
	for _, route := range oauth.Routes {
		var handler http.Handler
		if len(route.Middlewares) > 0 {
			n := negroni.New()
			for _, middleware := range route.Middlewares {
				n.Use(middleware)
			}
			n.Use(negroni.Wrap(route.HandlerFunc))
			handler = n
		} else {
			handler = route.HandlerFunc
		}

		subRouter.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Add routes for web pages
	subRouter = router.PathPrefix("/web").Subrouter()
	for _, route := range web.Routes {
		var handler http.Handler
		if len(route.Middlewares) > 0 {
			n := negroni.New()
			for _, middleware := range route.Middlewares {
				n.Use(middleware)
			}
			n.Use(negroni.Wrap(route.HandlerFunc))
			handler = n
		} else {
			handler = route.HandlerFunc
		}

		subRouter.Methods(route.Methods...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Set the router
	webApp.UseHandler(router)

	// Run the server on port 8080
	webApp.Run(":8080")
}
