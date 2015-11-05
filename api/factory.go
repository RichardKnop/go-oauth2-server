package api

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
)

// NewAPI - returns new *rest.Api with routes
func NewAPI(routes []*rest.Route) *rest.Api {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	return api
}
