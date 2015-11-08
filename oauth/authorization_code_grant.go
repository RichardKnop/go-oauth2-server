package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func authorizationCodeGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	code := r.FormValue("code")

	// Fetch an auth code from the database
	authCode := AuthCode{}
	if db.Where("code = ? AND client_id = ?", code, client.ID).First(&authCode).RecordNotFound() {
		api.Error(w, "Auth code not found", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := grantAccessToken(cnf, db, &authCode.Client, &authCode.User, authCode.Scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	respondWithAccessToken(w, cnf, accessToken, refreshToken)
}
