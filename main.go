package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/RichardKnop/go-oauth2-server/health"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/web"
	"github.com/areatech/go-fixtures"
	"github.com/codegangsta/cli"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	cliApp *cli.App
	app    *negroni.Negroni
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "area-api"
	cliApp.Usage = "Area Platform REST API"
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
	defer db.Close()

	// Disable logging
	db.LogMode(false)

	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) {
				migrate(db)
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) {
				loadData(c.Args(), cnf, db)
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) {
				runServer(cnf, db)
			},
		},
	}

	// Run the CLI app
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

func loadData(paths []string, cnf *config.Config, db *gorm.DB) {
	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		if err := fixtures.Load(data, db.DB(), cnf.Database.Type); err != nil {
			log.Fatal(err)
		}
	}
}

func runServer(cnf *config.Config, db *gorm.DB) {
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
