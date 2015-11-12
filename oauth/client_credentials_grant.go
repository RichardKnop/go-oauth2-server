package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) clientCredentialsGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	requestedScope := r.FormValue("scope")

	// Get the scope string
	scope, err := s.getScope(requestedScope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.GrantAccessToken(client, nil, scope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
