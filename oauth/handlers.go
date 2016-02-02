package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
)

var (
	errInvalidGrantType             = errors.New("Invalid grant type")
	errClientAuthenticationRequired = errors.New("Client authentication required")
)

// Handles all OAuth 2.0 grant types (POST /v1/oauth/tokens)
func (s *Service) tokensHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map of grant types against handler functions
	grantTypes := map[string]func(w http.ResponseWriter, r *http.Request, client *Client){
		"authorization_code": s.authorizationCodeGrant,
		"password":           s.passwordGrant,
		"client_credentials": s.clientCredentialsGrant,
		"refresh_token":      s.refreshTokenGrant,
	}

	// Check the grant type
	grantHandler, ok := grantTypes[r.Form.Get("grant_type")]
	if !ok {
		response.Error(w, errInvalidGrantType.Error(), http.StatusBadRequest)
		return
	}

	// Get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		response.UnauthorizedError(w, errClientAuthenticationRequired.Error())
		return
	}

	// Authenticate the client
	client, err := s.AuthClient(clientID, secret)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errClientAuthenticationRequired.Error())
		return
	}

	// Execute the correct function based on the grant type
	grantHandler(w, r, client)
}
