package web

import (
	"fmt"
	"net/http"
	"net/url"
)

// Redirects to a new path while keeping current request's query string
func redirectWithQueryString(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s", to, getQueryString(query)), http.StatusFound)
}

// Redirects to a new path with the query string moved to the URL fragment
func redirectWithFragment(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s#%s", to, query.Encode()), http.StatusFound)
}

// Returns string encoded query string of the request
func getQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}

// Helper function to handle redirecting failed or declined authorization
func errorRedirect(w http.ResponseWriter, r *http.Request, redirectURI *url.URL, err, state, responseType string) {
	query := redirectURI.Query()
	query.Set("error", err)
	if state != "" {
		query.Set("state", state)
	}
	if responseType == "code" {
		redirectWithQueryString(redirectURI.String(), query, w, r)
	}
	if responseType == "token" {
		redirectWithFragment(redirectURI.String(), query, w, r)
	}
}
