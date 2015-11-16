package web

import (
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the user session
	sessionService.ClearUserSession()

	// Redirect back to the login page
	redirectAndKeepQueryString("/web/login", w, r)
}
