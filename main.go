package main

import (
	"log"
	"net/http"

	"github.com/RichardKnop/go-microservice-example/migrations"
	"github.com/RichardKnop/go-microservice-example/service"
	"github.com/ant0ine/go-json-rest/rest"
)

func main() {
	log.Print("Starting up the Go Microservice Example")

	// Run database migrations
	migrations.RunAll()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/api/v1/users", service.UsersHandler),
		rest.Post("/api/v1/tokens", service.TokensHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
