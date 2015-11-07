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

	// Fetch refresh token from the database
	refreshToken := RefreshToken{}
	if db.Where(&RefreshToken{RefreshToken: token}).First(&refreshToken).RecordNotFound() {
		api.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	// Check refresh token hasn't expired
	if time.Now().After(refreshToken.ExpiresAt) {
		api.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	// Fetch the access token we are going to refresh
	accessToken := AccessToken{}
	if db.Where(&AccessToken{RefreshTokenID: refreshToken.ID}).Preload("Client").Preload("User").First(&accessToken).RecordNotFound() {
		api.Error(w, "Access token not found", http.StatusBadRequest)
		return
	}

	// Check the client IDs match
	if accessToken.Client.ClientID != client.ClientID {
		api.Error(w, "Client IDs mismatch", http.StatusBadRequest)
		return
	}

	// Requested scope CANNOT include any scope not originally granted
	if !scopeNotGreater(requestedScope, accessToken.Scope) {
		api.Error(w, "Invalid scope", http.StatusBadRequest)
		return
	}

	// Get the scope string
	scope, err := getScope(db, r.FormValue("scope"))
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	newAccessToken, err := grantAccessToken(cnf, db, &accessToken.Client, &accessToken.User, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	respondWithAccessToken(w, cnf, newAccessToken)
}
