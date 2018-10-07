package web

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/session"
)

// ErrIncorrectResponseType a form value for response_type was not set to token or code
var ErrIncorrectResponseType = errors.New("Response type not one of token or code")

func (s *Service) authorizeForm(w http.ResponseWriter, r *http.Request) {
	sessionService, client, _, responseType, _, err := s.authorizeCommon(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	query := r.URL.Query()
	query.Set("login_redirect_uri", r.URL.Path)
	renderTemplate(w, "authorize.html", map[string]interface{}{
		"error":       errMsg,
		"clientID":    client.Key,
		"queryString": getQueryString(query),
		"token":       responseType == "token",
	})
}

func (s *Service) authorize(w http.ResponseWriter, r *http.Request) {
	_, client, user, responseType, redirectURI, err := s.authorizeCommon(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the state parameter
	state := r.Form.Get("state")

	// Has the resource owner or authorization server denied the request?
	authorized := len(r.Form.Get("allow")) > 0
	if !authorized {
		errorRedirect(w, r, redirectURI, "access_denied", state, responseType)
		return
	}

	// Check the requested scope
	scope, err := s.oauthService.GetScope(r.Form.Get("scope"))
	if err != nil {
		errorRedirect(w, r, redirectURI, "invalid_scope", state, responseType)
		return
	}

	query := redirectURI.Query()

	// When response_type == "code", we will grant an authorization code
	if responseType == "code" {
		// Create a new authorization code
		authorizationCode, err := s.oauthService.GrantAuthorizationCode(
			client,                       // client
			user,                         // user
			s.cnf.Oauth.AuthCodeLifetime, // expires in
			redirectURI.String(),         // redirect URI
			scope,                        // scope
		)
		if err != nil {
			errorRedirect(w, r, redirectURI, "server_error", state, responseType)
			return
		}

		// Set query string params for the redirection URL
		query.Set("code", authorizationCode.Code)
		// Add state param if present (recommended)
		if state != "" {
			query.Set("state", state)
		}
		// And we're done here, redirect
		redirectWithQueryString(redirectURI.String(), query, w, r)
		return
	}

	// When response_type == "token", we will directly grant an access token
	if responseType == "token" {
		// Get access token lifetime from user input
		lifetime, err := strconv.Atoi(r.Form.Get("lifetime"))
		if err != nil {
			errorRedirect(w, r, redirectURI, "server_error", state, responseType)
			return
		}

		// Grant an access token
		accessToken, err := s.oauthService.GrantAccessToken(
			client,   // client
			user,     // user
			lifetime, // expires in
			scope,    // scope
		)
		if err != nil {
			errorRedirect(w, r, redirectURI, "server_error", state, responseType)
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
		redirectWithFragment(redirectURI.String(), query, w, r)
	}
}

func (s *Service) authorizeCommon(r *http.Request) (session.ServiceInterface, *models.OauthClient, *models.OauthUser, string, *url.URL, error) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		return nil, nil, nil, "", nil, err
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		return nil, nil, nil, "", nil, err
	}

	// Get the user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		return nil, nil, nil, "", nil, err
	}

	// Fetch the user
	user, err := s.oauthService.FindUserByUsername(
		userSession.Username,
	)
	if err != nil {
		return nil, nil, nil, "", nil, err
	}

	// Check the response_type is either "code" or "token"
	responseType := r.Form.Get("response_type")
	if responseType != "code" && responseType != "token" {
		return nil, nil, nil, "", nil, ErrIncorrectResponseType
	}

	// Fallback to the client redirect URI if not in query string
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = client.RedirectURI.String
	}

	// // Parse the redirect URL
	parsedRedirectURI, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return nil, nil, nil, "", nil, err
	}

	return sessionService, client, user, responseType, parsedRedirectURI, nil
}
