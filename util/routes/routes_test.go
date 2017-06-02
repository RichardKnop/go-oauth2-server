package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util/routes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

// helloWorldMiddleware is a test middleware that writes "hello world" to the
// response so we can check the middleware is registered
type helloWorldMiddleware struct{}

// ServeHTTP as per the negroni.Handler interface
func (m *helloWorldMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Write([]byte(fmt.Sprintf("hello world %s", id)))
	next(w, r)
}

func TestAddRoutes(t *testing.T) {
	var (
		router = mux.NewRouter()
		r      *http.Request
		err    error
		match  *mux.RouteMatch
		w      *httptest.ResponseRecorder
	)

	// Add a test GET route without a middleware
	routes.AddRoutes([]routes.Route{
		{
			Name:        "foobar_route",
			Method:      "GET",
			Pattern:     "/bar",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {},
		},
	}, router.PathPrefix("/foo").Subrouter())

	// Add a test PUT route with a middleware and a named parameter
	routes.AddRoutes([]routes.Route{
		{
			Name:        "helloworld_route",
			Method:      "PUT",
			Pattern:     "/world/{id:[0-9]+}",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {},
			Middlewares: []negroni.Handler{
				new(helloWorldMiddleware),
			},
		},
	}, router.PathPrefix("/hello").Subrouter())

	// Test the foobar_route
	r, err = http.NewRequest("GET", "http://1.2.3.4/foo/bar", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	// Test the route matches expected name
	match = new(mux.RouteMatch)
	router.Match(r, match)
	assert.Equal(t, "foobar_route", match.Route.GetName())

	// Test no middleware has been registered
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, "", w.Body.String())

	// Test the helloworld_route
	r, err = http.NewRequest("PUT", "http://1.2.3.4/hello/world/1", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	// Test the route matches expected name
	match = new(mux.RouteMatch)
	router.Match(r, match)
	assert.Equal(t, "helloworld_route", match.Route.GetName())

	// Test the helloWorldMiddleware has been registered
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	assert.Equal(t, "hello world 1", w.Body.String())
}

func TestRecoveryMiddlewareHandlesPanic(t *testing.T) {
	var (
		router = mux.NewRouter()
		app    = negroni.Classic()
		r      *http.Request
		err    error
		match  *mux.RouteMatch
	)

	// Add a test GET route without a middleware
	routes.AddRoutes([]routes.Route{
		{
			Name:    "panic_route",
			Method:  "GET",
			Pattern: "/panic",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				panic("oh no")
			},
		},
	}, router.PathPrefix("/foo").Subrouter())

	// Test the foobar_route
	r, err = http.NewRequest("GET", "http://1.2.3.4/foo/panic", nil)
	assert.NoError(t, err, "Request setup should not get an error")

	// Test the route matches expected name
	match = new(mux.RouteMatch)
	router.Match(r, match)
	assert.Equal(t, "panic_route", match.Route.GetName())

	// Test that panic does not crash the app
	app.UseHandler(router)
	app.ServeHTTP(httptest.NewRecorder(), r)
}
