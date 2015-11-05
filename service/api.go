package service

import (
	"log"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/ant0ine/go-json-rest/rest"
)

// NewAPI - returns new *rest.Api with routes
func NewAPI() *rest.Api {
	cnf := config.NewConfig()
	db, err := database.NewDatabase(cnf)
	if err != nil {
		log.Fatal(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/api/v1/users", func(w rest.ResponseWriter, r *rest.Request) {
			registerUser(w, r, cnf, db)
		}),
		rest.Post("/api/v1/tokens", func(w rest.ResponseWriter, r *rest.Request) {
			tokensHandler(w, r, cnf, db)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api
}
