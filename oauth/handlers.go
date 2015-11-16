package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

// Map of grant types against handler functions
var grantTypes = map[string]func(w http.ResponseWriter, r *http.Request, client *Client){
	"authorization_code": theService.authorizationCodeGrant,
	"password":           theService.passwordGrant,
	"client_credentials": theService.clientCredentialsGrant,
	"refresh_token":      theService.refreshTokenGrant,
}

// Handles all OAuth 2.0 grant types (POST /oauth2/api/v1/tokens
func handleTokens(w http.ResponseWriter, r *http.Request) {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		json.Error(w, "Parse form error", http.StatusInternalServerError)
		return
	}

	// Check the grant type
	grantHandler, ok := grantTypes[r.Form.Get("grant_type")]
	if !ok {
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

	// Execute the correct function based on the grant type
	grantHandler(w, r, client)
}
