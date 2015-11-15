package web

import (
	"log"
	"net/http"
)

func loginForm(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := newSessionService(theService.cnf, r, w)
	if err := sessionService.initSession("user_session"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "login.tmpl", map[string]interface{}{
		"error": sessionService.getFlashMessage(),
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := newSessionService(theService.cnf, r, w)
	if err := sessionService.initSession("user_session"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the submitted form data
	r.ParseForm()
	username := r.Form["email"][0]
	password := r.Form["password"][0]

	// Fetch the trusted client
	client, err := theService.oauthService.FindClientByClientID(
		theService.cnf.TrustedClient.ClientID,
	)
	if err != nil {
		sessionService.setFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Authenticate the user
	user, err := theService.oauthService.AuthUser(username, password)
	if err != nil {
		sessionService.setFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Default scope
	scope := "read_write"

	// Grant an access token
	accessToken, err := theService.oauthService.GrantAccessToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.setFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// Get a refresh token
	refreshToken, err := theService.oauthService.GetOrCreateRefreshToken(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.setFlashMessage(err.Error())
		http.Redirect(w, r, "/web/login", http.StatusFound)
		return
	}

	// TODO - store access and refresh token in the session

	log.Print(accessToken)
	log.Print(refreshToken)
	// http.Redirect(w, r, "/web/authorize", http.StatusFound)
}
