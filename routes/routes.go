package routes

import (
	"net/http"

	"github.com/codegangsta/negroni"
)

// Route ...
type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []negroni.Handler
}
