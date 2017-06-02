package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util/response"
	"github.com/RichardKnop/jsonhal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// TestGetErrorExpectedResponse ...
func TestGetErrorExpectedResponse(t *testing.T, router *mux.Router, url, accessToken, msg string, code int, assertExpectations func()) {
	TestErrorExpectedResponse(t, router, "GET", url, nil, accessToken, msg, code, assertExpectations)
}

// TestPutErrorExpectedResponse ...
func TestPutErrorExpectedResponse(t *testing.T, router *mux.Router, url string, data io.Reader, accessToken, msg string, code int, assertExpectations func()) {
	TestErrorExpectedResponse(t, router, "PUT", url, data, accessToken, msg, code, assertExpectations)
}

// TestPostErrorExpectedResponse ...
func TestPostErrorExpectedResponse(t *testing.T, router *mux.Router, url string, data io.Reader, accessToken, msg string, code int, assertExpectations func()) {
	TestErrorExpectedResponse(t, router, "POST", url, data, accessToken, msg, code, assertExpectations)
}

// TestErrorExpectedResponse is the generic test code for testing for a bad response
func TestErrorExpectedResponse(t *testing.T, router *mux.Router, method, url string, data io.Reader, accessToken, msg string, code int, assertExpectations func()) {
	// Prepare a request
	r, err := http.NewRequest(
		method,
		url,
		data,
	)
	assert.NoError(t, err)

	// Optionally add a bearer token to headers
	if accessToken != "" {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	// And serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	TestResponseForError(t, w, msg, code)

	assertExpectations()
}

// TestResponseForError tests a response w to see if it returned an error msg with http code
func TestResponseForError(t *testing.T, w *httptest.ResponseRecorder, msg string, code int) {
	if code != w.Code {
		log.Print(w.Body.String())
	}
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	assert.NotNil(t, w)
	TestResponseBody(t, w, getErrorJSON(msg))
}

// TestEmptyResponse tests an empty 204 response
func TestEmptyResponse(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, 204, w.Code)
	TestResponseBody(t, w, "")
}

// TestResponseObject tests response body is equal to expected object in JSON form
func TestResponseObject(t *testing.T, w *httptest.ResponseRecorder, expected interface{}, code int) {
	if code != w.Code {
		log.Print(w.Body.String())
	}
	assert.Equal(
		t,
		code,
		w.Code,
		fmt.Sprintf("Expected a %d response but got %d", code, w.Code),
	)
	jsonBytes, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.NotNil(t, w)
	assert.Equal(
		t,
		string(jsonBytes),
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

// TestResponseBody tests response body is equal to expected string
func TestResponseBody(t *testing.T, w *httptest.ResponseRecorder, expected string) {
	assert.Equal(
		t,
		expected,
		strings.TrimRight(w.Body.String(), "\n"),
		"Should have returned correct body text",
	)
}

// TestListValidResponse ...
func TestListValidResponse(t *testing.T, router *mux.Router, path, entity, accessToken string, items []interface{}, assertExpectations func()) {
	TestListValidResponseWithParams(t, router, path, entity, accessToken, items, assertExpectations, nil)
}

// TestListValidResponseWithParams tests a list endpoint for a valid response with default settings
func TestListValidResponseWithParams(t *testing.T, router *mux.Router, path, entity, accessToken string, items []interface{}, assertExpectations func(), params map[string]string) {
	u, err := url.Parse(fmt.Sprintf("http://1.2.3.4/v1/%s", path))
	assert.NoError(t, err)

	// add any params
	for k, v := range params {
		q := u.Query()
		q.Set(k, v)
		u.RawQuery = q.Encode()
	}

	// Prepare a request
	r, err := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)
	assert.NoError(t, err)

	if accessToken != "" {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	// And serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Check that the mock object expectations were met
	assertExpectations()

	// Check the status code
	assert.Equal(t, http.StatusOK, w.Code)

	baseURI := u.RequestURI()

	q := u.Query()
	q.Set("page", "1")
	u.RawQuery = q.Encode()

	pagedURI := u.RequestURI()

	expected := &response.ListResponse{
		Hal: jsonhal.Hal{
			Links: map[string]*jsonhal.Link{
				"self": {
					Href: baseURI,
				},
				"first": {
					Href: pagedURI,
				},
				"last": {
					Href: pagedURI,
				},
				"prev": new(jsonhal.Link),
				"next": new(jsonhal.Link),
			},
			Embedded: map[string]jsonhal.Embedded{
				entity: jsonhal.Embedded(items),
			},
		},
		Count: uint(len(items)),
		Page:  1,
	}
	expectedJSON, err := json.Marshal(expected)

	if assert.NoError(t, err, "JSON marshalling failed") {
		TestResponseBody(t, w, string(expectedJSON))
	}
}

func getErrorJSON(msg string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", msg)
}
