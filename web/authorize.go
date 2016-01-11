package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (s *Service) authorizeForm(w http.ResponseWriter, r *http.Request) {
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
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "authorize.html", map[string]interface{}{
		"error":       errMsg,
		"clientID":    client.Key,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) authorize(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the user
	user, err := s.oauthService.FindUserByUsername(
		userSession.Username,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check the response_type is either "code" or "token"
	responseType := r.Form.Get("response_type")
	if responseType != "code" && responseType != "token" {
		http.Error(w, "Invalid response type", http.StatusBadRequest)
		return
	}

	// Fallback to the client redirect URI if not in query string
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		value, err := client.RedirectURI.Value()
		if err == nil {
			clientRedirectURI, ok := value.(string)
			if ok {
				redirectURI = clientRedirectURI
			}
		}
	}
	// // Parse the redirect URL
	parsedRedirectURI, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the state parameter
	state := r.Form.Get("state")

	// Has the resource owner or authorization server denied the request?
	authorized := len(r.Form.Get("authorize")) > 0
	if !authorized {
		errorRedirect(w, r, parsedRedirectURI, "access_denied", state, responseType)
		return
	}

	// Check the requested scope
	scope, err := s.oauthService.GetScope(r.Form.Get("scope"))
	if err != nil {
		errorRedirect(w, r, parsedRedirectURI, "invalid_scope", state, responseType)
		return
	}

	query := parsedRedirectURI.Query()

	// When response_type == "code", we will grant an authorization code
	if responseType == "code" {
		// Create a new authorization code
		authorizationCode, err := s.oauthService.GrantAuthorizationCode(
			client, // client
			user,   // user
			r.Form.Get("redirect_uri"), // redirect URI
			scope, // scope
		)
		if err != nil {
			log.Print(err)
			errorRedirect(w, r, parsedRedirectURI, "server_error", state, responseType)
			return
		}

		// Set query string params for the redirection URL
		query.Set("code", authorizationCode.Code)
		// Add state param if present (recommended)
		if state != "" {
			query.Set("state", state)
		}
		// And we're done here, redirect
		redirectWithQueryString(parsedRedirectURI.String(), query, w, r)
		return
	}

	// When response_type == "token", we will directly grant an access token
	if responseType == "token" {
		// Grant an access token
		accessToken, err := s.oauthService.GrantAccessToken(
			client, // client
			user,   // user
			scope,  // scope
		)
		if err != nil {
			log.Print(err)
			errorRedirect(w, r, parsedRedirectURI, "server_error", state, responseType)
			return
		}

		// Set query string params for the redirection URL
		query.Set("access_token", accessToken.Token)
		query.Set("expires_in", fmt.Sprintf("%d", s.cnf.Oauth.AccessTokenLifetime))
		query.Set("token_type", "Bearer")
		query.Set("scope", scope)
		// Add state param if present (recommended)
		if state != "" {
			query.Set("state", state)
		}
		// And we're done here, redirect
		redirectWithFragment(parsedRedirectURI.String(), query, w, r)
	}
}
