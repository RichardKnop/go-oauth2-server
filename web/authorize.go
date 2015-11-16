package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
		"queryString": getQueryString(r.URL.Query()),
	})
}

func authorize(w http.ResponseWriter, r *http.Request) {
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

	// The resource owner or authorization server denied the request
	declined := strings.ToLower(r.Form.Get("authorize")) != "authorize"
	if declined {
		authorizeErrorRedirect(w, r, parsedRedirectURI, "access_denied", state)
		return
	}

	// Check the response_type is either "code" or "token"
	responseType := r.Form.Get("response_type")
	if responseType != "code" && responseType != "token" {
		authorizeErrorRedirect(w, r, parsedRedirectURI, "unsupported_response_type", state)
		return
	}

	// Check the requested scope
	scope, err := theService.oauthService.GetScope(r.Form.Get("scope"))
	if err != nil {
		authorizeErrorRedirect(w, r, parsedRedirectURI, "invalid_scope", state)
		return
	}

	// Get the user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	query := parsedRedirectURI.Query()

	// When response_type == "code", we will grant an authorization code
	if responseType == "code" {
		// Create a new authorization code
		authorizationCode, err := theService.oauthService.GrantAuthorizationCode(
			client,
			user,
			parsedRedirectURI.String(),
			scope,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set query string params for the redirection URL
		query.Set("code", authorizationCode.Code)
	}

	// When response_type == "token", we will directly grant an access token
	if responseType == "token" {
		// Grant an access token
		accessToken, err := theService.oauthService.GrantAccessToken(
			client,
			user,
			scope,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get a refresh token
		refreshToken, err := theService.oauthService.GetOrCreateRefreshToken(
			client,
			user,
			scope,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set query string params for the redirection URL
		query.Set("access_token", accessToken.Token)
		query.Set("expires_in", fmt.Sprintf("%d", theService.cnf.Oauth.AccessTokenLifetime))
		query.Set("token_type", "Bearer")
		query.Set("refresh_token", refreshToken.Token)
	}

	// Add state param if present (recommended)
	if state != "" {
		query.Set("state", state)
	}
	// And we're done here, redirect
	redirectWithQueryString(parsedRedirectURI.String(), query, w, r)
}
