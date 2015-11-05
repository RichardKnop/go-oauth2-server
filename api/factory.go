package api

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
)

// NewAPI - returns new *rest.Api with routes
func NewAPI(stack []rest.Middleware, routes []*rest.Route) *rest.Api {
	api := rest.NewApi()
	api.Use(stack...)
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api
}
