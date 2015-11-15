package web

import (
	"log"
	"net/http"
)

func loginForm(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "login.tmpl", map[string]interface{}{
		"error": sessionService.getFlashMessage(r, w),
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the submitted form data
	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	// Fetch the trusted client
	client, err := s.oauthService.FindClientByClientID(s.cnf.TrustedClient.ClientID)
	if err != nil {
		sessionService.setFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Authenticate the user
	user, err := s.oauthService.AuthUser(username, password)
	if err != nil {
		sessionService.setFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Default scope
	scope := "read_write"

	// Grant an access token
	accessToken, err := s.oauthService.GrantAccessToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.setFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Get a refresh token
	refreshToken, err := s.oauthService.GetOrCreateRefreshToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.setFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// TODO - store access and refresh token in the session

	log.Print(accessToken)
	log.Print(refreshToken)
	// http.Redirect(w, r, "/web/authorize", http.StatusFound)
}
