package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RichardKnop/go-microservice-example/migrations"
	"github.com/RichardKnop/go-microservice-example/service"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-microservice-example"
	app.Usage = "OAuth 2.0 Go microservice"
	app.Author = "Richard Knop"
	app.Email = "risoknop@gmail.com"
	app.Version = "0.0.0"

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
				api := service.NewAPI()
				log.Print("Listening on port 8080")
				log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
			},
		},
	}

	app.Run(os.Args)
}
