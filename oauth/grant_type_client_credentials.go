package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) clientCredentialsGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		client,    // client
		new(User), // empty user
		scope,     // scope
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		client,    // client
		new(User), // empty user
		scope,     // scope
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
