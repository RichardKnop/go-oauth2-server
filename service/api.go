package service

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
)

// NewAPI - returns new *rest.Api with routes
func NewAPI() *rest.Api {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Post("/api/v1/users", registerUser),
		rest.Post("/api/v1/tokens", tokensHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api
}
