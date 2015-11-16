package web

import (
	"fmt"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
)

func noLoginRequired(w http.ResponseWriter, r *http.Request) *session.Service {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return sessionService
}

func loginRequired(w http.ResponseWriter, r *http.Request) *session.Service {
	sessionService := noLoginRequired(w, r)
	if sessionService == nil {
		return nil
	}

	// If the user is not logged in, redirect to the login page
	if err := sessionService.IsLoggedIn(); err != nil {
		redirectAndKeepQueryString("/web/login", w, r)
		return nil
	}

	return sessionService
}

func redirectAndKeepQueryString(path string, w http.ResponseWriter, r *http.Request) {
	redirectURL := path
	if len(r.URL.Query()) > 0 {
		redirectURL = fmt.Sprintf("%s?%s", redirectURL, r.URL.Query().Encode())
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
