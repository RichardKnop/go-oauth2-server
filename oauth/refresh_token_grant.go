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

	refreshToken := RefreshToken{}
	if db.Where(&RefreshToken{RefreshToken: token}).First(&refreshToken).RecordNotFound() {
		api.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		api.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	oldAccessToken := AccessToken{}
	if db.Where(&AccessToken{ClientID: client.ID, RefreshTokenID: refreshToken.ID}).Preload("Client").Preload("User").First(&oldAccessToken).RecordNotFound() {
		api.Error(w, "Access token not found", http.StatusBadRequest)
		return
	}

	// Requested scope CANNOT include any scope not originally granted
	if !scopeNotGreater(r.FormValue("scope"), oldAccessToken.Scope) {
		api.Error(w, "Invalid scope", http.StatusBadRequest)
		return
	}

	scope, err := getScope(db, r.FormValue("scope"))
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newAccessToken, err := grantAccessToken(cnf, db, &oldAccessToken.Client, &oldAccessToken.User, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	respondWithAccessToken(w, cnf, newAccessToken)
}
