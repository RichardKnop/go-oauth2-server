package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RichardKnop/go-microservice-example/migrations"
	"github.com/RichardKnop/go-microservice-example/service"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-microservice-example"
	app.Usage = "OAuth 2.0 Go microservice"

	app.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) {
				migrations.RunAll()
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) {
				api := rest.NewApi()
				api.Use(rest.DefaultDevStack...)
				router, err := rest.MakeRouter(
					rest.Post("/api/v1/users", service.RegisterUser),
					rest.Post("/api/v1/tokens", service.TokensHandler),
				)
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
