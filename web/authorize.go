package web

import (
	"log"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
)

func authorizeForm(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the user is not logged in, redirect to the login page
	if err := sessionService.IsLoggedIn(); err != nil {
		http.Redirect(w, r, "/web/login?"+r.URL.Query().Encode(), http.StatusFound)
		return
	}

	// If there is a flash message, just render the template with the error
	if err := sessionService.GetFlashMessage(); err != nil {
		renderTemplate(w, "authorize.tmpl", map[string]interface{}{"error": err})
		return
	}

	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check the response_type is either "code" or "token"
	responseType := r.Form.Get("response_type")
	if responseType != "code" && responseType != "token" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	// Fetch the client
	client, err := theService.oauthService.FindClientByClientID(
		r.Form.Get("client_id"),
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Render the template
	renderTemplate(w, "authorize.tmpl", map[string]interface{}{
		"clientID": client.ClientID,
	})
}

func authorize(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the user is not logged in, redirect to the login page
	if err := sessionService.IsLoggedIn(); err != nil {
		http.Redirect(w, r, "/web/login?"+r.URL.Query().Encode(), http.StatusFound)
		return
	}

	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get required parameters from the query string
	responseType := r.Form.Get("response_type")
	clientID := r.Form.Get("client_id")
	redirectURI := r.Form.Get("redirect_uri")
	scope := r.Form.Get("scope")
	state := r.Form.Get("state")

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
	client, err := theService.oauthService.FindClientByClientID(clientID)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	log.Print(client)

	// TODO
}
