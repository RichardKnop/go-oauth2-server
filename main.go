package main

import (
	"log"
	"os"

	"github.com/adam-hanna/go-oauth2-server/cmd"
	"github.com/urfave/cli"
)

var (
	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "go-oauth2-server"
	cliApp.Usage = "Go OAuth 2.0 Server"
	cliApp.Author = "Richard Knop"
	cliApp.Email = "risoknop@gmail.com"
	cliApp.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "run migrations",
			Action: func(c *cli.Context) error {
				return cmd.Migrate()
			},
		},
		{
			Name:  "loaddata",
			Usage: "load data from fixture",
			Action: func(c *cli.Context) error {
				return cmd.LoadData(c.Args())
			},
		},
		{
			Name:  "runserver",
			Usage: "run web server",
			Action: func(c *cli.Context) error {
				return cmd.RunServer()
			},
		},
	}

	// Run the CLI app
	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
