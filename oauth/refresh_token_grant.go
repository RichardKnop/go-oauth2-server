package oauth

import (
	"net/http"
	"time"

	"github.com/RichardKnop/go-oauth2-server/errors"
	"github.com/ant0ine/go-json-rest/rest"
)

func (s *service) refreshTokenGrant(w rest.ResponseWriter, r *rest.Request, client *Client) {
	token := r.FormValue("refresh_token")
	requestedScope := r.FormValue("scope")

	// Fetch a refresh token from the database
	theRefreshToken := RefreshToken{}
	if s.db.Where(&RefreshToken{
		Token:    token,
		ClientID: clientIDOrNull(client),
	}).Preload("Client").Preload("User").First(&theRefreshToken).RecordNotFound() {
		errors.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	// Check the refresh token hasn't expired
	if time.Now().After(theRefreshToken.ExpiresAt) {
		errors.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	// Requested scope CANNOT include any scope not originally granted
	if !s.scopeNotGreater(requestedScope, theRefreshToken.Scope) {
		errors.Error(w, "Invalid scope", http.StatusBadRequest)
		return
	}

	// Get the scope string
	scope, err := s.getScope(requestedScope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(&theRefreshToken.Client, &theRefreshToken.User, scope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeAccessToken(w, s.cnf.AccessTokenLifetime, accessToken, refreshToken)
}
