package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

// Handles all OAuth 2.0 grant types (POST /oauth2/api/v1/tokens
func handleTokens(w http.ResponseWriter, r *http.Request) {
	// Check the grant type
	grantTypes := map[string]bool{
		"authorization_code": true,
		"password":           true,
		"client_credentials": true,
		"refresh_token":      true,
	}
	if !grantTypes[r.FormValue("grant_type")] {
		json.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	// Authenticate the client
	client, err := s.authClient(r)
	if err != nil {
		json.UnauthorizedError(w, err.Error())
		return
	}

	grants := map[string]func(){
		"authorization_code": func() { s.authorizationCodeGrant(w, r, client) },
		"password":           func() { s.passwordGrant(w, r, client) },
		"client_credentials": func() { s.clientCredentialsGrant(w, r, client) },
		"refresh_token":      func() { s.refreshTokenGrant(w, r, client) },
	}
	grants[r.FormValue("grant_type")]()
}
