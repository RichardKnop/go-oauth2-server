package web

import (
	"net/http"
	"net/url"
)

func authorizeForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	renderTemplate(w, "authorize.html", map[string]interface{}{
		"error":       sessionService.GetFlashMessage(),
		"clientID":    client.ClientID,
		"queryString": getQueryString(r),
	})
}

func authorize(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get required parameters from the query string
	responseType := r.Form.Get("response_type")
	scope := r.Form.Get("scope")
	state := r.Form.Get("state")

	// Check the response_type is either "code" or "token"
	if responseType != "code" && responseType != "token" {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	redirectURI, err := url.Parse(r.Form.Get("redirect_uri"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare the client ID from session with the one from the query string
	if userSession.ClientID != client.ClientID {
		http.Error(w, "Client ID mismatch", http.StatusInternalServerError)
		return
	}

	// Fetch the user
	user, err := theService.oauthService.FindUserByUsername(
		userSession.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if responseType == "code" {
		// Create a new authorization code
		authorizationCode, err := theService.oauthService.GrantAuthorizationCode(
			client,
			user,
			scope,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		redirectURI.Query().Add("code", authorizationCode.Code)

	}

	if responseType == "implicit" {
		redirectURI.Query().Add("access_token", userSession.AccessToken)
		redirectURI.Query().Add("expires_in", string(theService.cnf.Oauth.AccessTokenLifetime))
		redirectURI.Query().Add("token_type", "Bearer")
		redirectURI.Query().Add("refresh_token", userSession.RefreshToken)
	}

	redirectURI.Query().Add("state", state)
	http.Redirect(w, r, redirectURI.RequestURI(), 302)
}
