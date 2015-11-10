package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/errors"
	"github.com/ant0ine/go-json-rest/rest"
)

func (s *service) passwordGrant(w rest.ResponseWriter, r *rest.Request, client *Client) {
	requestedScope := r.FormValue("scope")

	// Authenticate the user
	user, err := s.authUser(r.Request)
	if err != nil {
		errors.UnauthorizedError(w, err.Error())
		return
	}

	// Get the scope string
	scope, err := s.getScope(requestedScope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(client, user, scope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeAccessToken(w, s.cnf.AccessTokenLifetime, accessToken, refreshToken)
}
