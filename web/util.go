package web

import (
	"fmt"
	"net/http"
)

// Redirects to a new path while keeping current request's query string
func redirectAndKeepQueryString(path string, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s", path, getQueryString(r)), http.StatusFound)
}

// Returns string encoded query string of the request
func getQueryString(r *http.Request) string {
	if len(r.URL.Query()) > 0 {
		return fmt.Sprintf("?%s", r.URL.Query().Encode())
	}
	return ""
}
