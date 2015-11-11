package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	requestedScope := r.FormValue("scope")

	// Get user credentials from from the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authenticate the user
	user, err := s.AuthUser(username, password)
	if err != nil {
		json.UnauthorizedError(w, err.Error())
		return
	}

	// Get the scope string
	scope, err := s.getScope(requestedScope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(client, user, scope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
