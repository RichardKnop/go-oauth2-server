package web

import (
	"log"
	"net/http"
)

func authorizeForm(w http.ResponseWriter, r *http.Request) {
	session, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseType := r.URL.Query().Get("response_type")
	if responseType != "code" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	log.Print(clientID)
	log.Print(redirectURI)
	log.Print(scope)
	log.Print(state)

	data := map[string]interface{}{}
	data["error"] = getLastFlashMessage(session, r, w)
	data["client_id"] = "test_client"

	renderTemplate(w, "authorize.tmpl", data)
}

func authorize(w http.ResponseWriter, r *http.Request) {
	_, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseType := r.URL.Query().Get("response_type")
	if responseType != "code" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	scope := r.URL.Query().Get("scope")
	state := r.URL.Query().Get("state")

	log.Print(clientID)
	log.Print(redirectURI)
	log.Print(scope)
	log.Print(state)
}
