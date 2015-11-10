package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/errors"
	"github.com/ant0ine/go-json-rest/rest"
)

func (s *service) clientCredentialsGrant(w rest.ResponseWriter, r *rest.Request, client *Client) {
	requestedScope := r.FormValue("scope")

	// Get the scope string
	scope, err := s.getScope(requestedScope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(client, nil, scope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeAccessToken(w, s.cnf.AccessTokenLifetime, accessToken, refreshToken)
}
