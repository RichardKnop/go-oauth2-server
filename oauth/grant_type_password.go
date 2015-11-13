package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Double check the grant type
	if r.FormValue("grant_type") != "password" {
		json.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	// Get user credentials from from the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authenticate the user
	user, err := s.AuthUser(username, password)
	if err != nil {
		// For security reasons, return a general error message
		json.UnauthorizedError(w, "User authentication required")
		return
	}

	// Get the scope string
	scope, err := s.getScope(r.FormValue("scope"))
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(client, user, scope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(client, user, scope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
