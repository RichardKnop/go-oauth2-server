package routes

import (
	"log"
	"net/http"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// helloWorldMiddleware is a test middleware that does nothing
type helloWorldMiddleware struct{}

// ServeHTTP as per the negroni.Handler interface
func (m *helloWorldMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func TestAddRoutes(t *testing.T) {
	router := mux.NewRouter().StrictSlash(true)

	// Add a test GET route without a middleware
	AddRoutes([]Route{
		Route{
			Name:        "foobar_route",
			Methods:     []string{"GET"},
			Pattern:     "/bar",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {},
		},
	}, router.PathPrefix("/foo").Subrouter())

	// Add a test POST route with a middleware
	AddRoutes([]Route{
		Route{
			Name:        "helloworld_route",
			Methods:     []string{"POST"},
			Pattern:     "/world",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {},
			Middlewares: []negroni.Handler{
				new(helloWorldMiddleware),
			},
		},
	}, router.PathPrefix("/hello").Subrouter())

	var match *mux.RouteMatch
	var r *http.Request
	var err error

	// Test the foobar_route
	r, err = http.NewRequest("GET", "http://1.2.3.4/foo/bar", nil)
	if err != nil {
		log.Fatal(err)
	}
	match = new(mux.RouteMatch)
	router.Match(r, match)
	assert.Equal(t, "foobar_route", match.Route.GetName())

	// Test the helloworld_route
	r, err = http.NewRequest("POST", "http://1.2.3.4/hello/world", nil)
	if err != nil {
		log.Fatal(err)
	}
	match = new(mux.RouteMatch)
	router.Match(r, match)
	assert.Equal(t, "helloworld_route", match.Route.GetName())
	// TODO - test the helloWorldMiddleware was added to the route
}
