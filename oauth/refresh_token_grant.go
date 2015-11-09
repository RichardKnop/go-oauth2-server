package oauth

import (
	"net/http"
	"time"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func refreshTokenGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	token := r.FormValue("refresh_token")
	requestedScope := r.FormValue("scope")

	// Fetch a refresh token from the database
	theRefreshToken := RefreshToken{}
	if db.Where(&RefreshToken{
		Token:    token,
		ClientID: clientIDOrNull(client),
	}).Preload("Client").Preload("User").First(&theRefreshToken).RecordNotFound() {
		api.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	// Check the refresh token hasn't expired
	if time.Now().After(theRefreshToken.ExpiresAt) {
		api.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	// Requested scope CANNOT include any scope not originally granted
	if !scopeNotGreater(requestedScope, theRefreshToken.Scope) {
		api.Error(w, "Invalid scope", http.StatusBadRequest)
		return
	}

	// Get the scope string
	scope, err := getScope(db, requestedScope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := grantAccessToken(cnf, db, &theRefreshToken.Client, &theRefreshToken.User, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	respondWithAccessToken(w, cnf, accessToken, refreshToken)
}
