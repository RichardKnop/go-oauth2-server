package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/database"
	"github.com/RichardKnop/go-oauth2-server/migrate"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-oauth2-server"
	app.Usage = "OAuth 2.0 Server"
	app.Author = "Richard Knop"
	app.Email = "risoknop@gmail.com"
	app.Version = "0.0.0"

	cnf := config.NewConfig()
	db, err := database.NewDatabase(cnf)
	if err != nil {
		log.Fatal(err)
	}

	app.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) {
				if err := migrate.Bootstrap(db); err != nil {
					log.Fatal(err)
				}
				if err := oauth.MigrateAll(db); err != nil {
					log.Fatal(err)
				}
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) {
				api := rest.NewApi()
				api.Use([]rest.Middleware{
					&rest.AccessLogApacheMiddleware{
						Format: rest.CombinedLogFormat,
					},
					&rest.TimerMiddleware{},
					&rest.RecorderMiddleware{},
					&rest.PoweredByMiddleware{},
					&rest.RecoverMiddleware{},
					&rest.GzipMiddleware{},
				}...)
				router, err := rest.MakeRouter(oauth.NewRoutes(cnf, db)...)
				if err != nil {
					log.Fatal(err)
				}
				api.SetApp(router)
				log.Print("Listening on port 8080")
				log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
			},
		},
	}

	app.Run(os.Args)
}
