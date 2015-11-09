package oauth

import (
	"log"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func implicitGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	code := r.FormValue("code")

	// Fetch an auth code from the database
	authCode := AuthCode{}
	if db.Where(&AuthCode{
		Code:     code,
		ClientID: clientIDOrNull(client),
	}).First(&authCode).RecordNotFound() {
		api.Error(w, "Auth code not found", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := grantAccessToken(cnf, db, &authCode.Client, &authCode.User, authCode.Scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print(accessToken)
	log.Print(refreshToken)

	// TODO redirect
}
