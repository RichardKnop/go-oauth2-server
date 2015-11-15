package web

import (
	"log"
	"net/http"
)

func authorizeForm(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get required parameters from the query string
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	log.Print(responseType)
	log.Print(clientID)
	log.Print(redirectURI)
	log.Print(scope)
	log.Print(state)

	// Check the response_type is either "code" or "token"
	if responseType != "code" && responseType != "token" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	// Fetch the client
	client, err := s.oauthService.FindClientByClientID(clientID)
	if err != nil {
		sessionService.addFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/authorize", http.StatusFound)
		return
	}

	// Render the template
	renderTemplate(w, "authorize.tmpl", map[string]interface{}{
		"error":  sessionService.getLastFlashMessage(r, w),
		"client": client,
	})
}

func authorize(w http.ResponseWriter, r *http.Request) {
	// Initialise a new session service
	sessionService := newSessionService(s.cnf)
	if err := sessionService.initSession(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get required parameters from the query string
	responseType := r.URL.Query().Get("response_type")
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	log.Print(responseType)
	log.Print(clientID)
	log.Print(redirectURI)
	log.Print(scope)
	log.Print(state)

	// Check the response_type is either "code" or "token"
	if responseType != "code" && responseType != "token" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	// Fetch the client
	client, err := s.oauthService.FindClientByClientID(clientID)
	if err != nil {
		sessionService.addFlashMessage(err.Error(), r, w)
		http.Redirect(w, r, "/web/authorize", http.StatusFound)
		return
	}

	log.Print(client)

	// TODO
}
