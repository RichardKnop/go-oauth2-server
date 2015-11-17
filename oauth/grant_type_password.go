package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Authenticate the user
	user, err := s.AuthUser(r.Form.Get("username"), r.Form.Get("password"))
	if err != nil {
		// For security reasons, return a general error message
		json.UnauthorizedError(w, "User authentication required")
		return
	}

	// Get the scope string
	scope, err := s.GetScope(r.Form["scope"][0])
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		client, // client
		user,   // user
		scope,  // scope
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		client,
		user,
		scope,
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
