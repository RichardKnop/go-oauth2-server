package web

import (
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {
	sessionService := loginRequired(w, r)
	if sessionService == nil {
		return
	}

	// Logout the user
	sessionService.LogOut()

	// Redirect back to the login page
	redirectAndKeepQueryString("/web/login", w, r)
}
