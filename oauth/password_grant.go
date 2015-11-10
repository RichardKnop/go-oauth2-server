package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	requestedScope := r.FormValue("scope")

	// Authenticate the user
	user, err := s.authUser(r)
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
