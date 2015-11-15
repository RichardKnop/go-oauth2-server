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

	// Get client credentials from basic auth
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		json.UnauthorizedError(w, "Client authentication required")
		return
	}

	// Authenticate the client
	client, err := theService.AuthClient(clientID, clientSecret)
	if err != nil {
		// For security reasons, return a general error message
		json.UnauthorizedError(w, "Client authentication required")
		return
	}

	// Map of grant types against handler functions
	grants := map[string]func(){
		"authorization_code": func() {
			theService.authorizationCodeGrant(w, r, client)
		},
		"password": func() {
			theService.passwordGrant(w, r, client)
		},
		"client_credentials": func() {
			theService.clientCredentialsGrant(w, r, client)
		},
		"refresh_token": func() {
			theService.refreshTokenGrant(w, r, client)
		},
	}

	// Execute the correct function based on the grant type
	grants[r.FormValue("grant_type")]()
}
