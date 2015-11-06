package oauth2

import (
	"net/http"
	"strings"
	"time"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Handles all OAuth 2.0 grant types
func tokens(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	// Check grant type
	if !checkGrantType(r) {
		api.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	// Client authentication required
	client, err := authClient(r, db)
	if err != nil {
		api.UnauthorizedError(w, err.Error())
		return
	}

	grants := map[string]func(){
		"password":           func() { password(w, r, cnf, db, client) },
		"client_credentials": func() { clientCredentials(w, r, cnf, db, client) },
		"refresh_token":      func() { refreshToken(w, r, cnf, db, client) },
	}
	grants[r.FormValue("grant_type")]()
}

// Checks grant type parameter from posted form data
func checkGrantType(r *rest.Request) bool {
	grantTypes := map[string]bool{
		// "authorization_code": true,
		// "implicit":           true,
		"password":           true,
		"client_credentials": true,
		"refresh_token":      true,
	}
	return grantTypes[r.FormValue("grant_type")]
}

// Grants user credentials access token
func password(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	// User authentication required
	user, err := authUser(r, db)
	if err != nil {
		api.UnauthorizedError(w, err.Error())
		return
	}

	grantAccessToken(w, cnf, db, client.ID, user.ID)
}

// Grants client credentials access token
func clientCredentials(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	grantAccessToken(w, cnf, db, client.ID, -1)
}

// Refreshes access token
func refreshToken(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	token := r.FormValue("refresh_token")

	refreshToken := RefreshToken{}
	if db.Where(&RefreshToken{RefreshToken: token}).First(&refreshToken).RecordNotFound() {
		api.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	if refreshToken.ExpiresAt.After(time.Now()) {
		api.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	accessToken := AccessToken{}
	if db.Where(&AccessToken{RefreshTokenID: refreshToken.ID}).First(&accessToken).RecordNotFound() {
		api.Error(w, "Access token not found", http.StatusBadGateway)
		return
	}

	// Delete old access / refresh token
	db.Delete(&refreshToken)
	db.Delete(&accessToken)

	grantAccessToken(w, cnf, db, accessToken.ClientID, accessToken.UserID)
}

// Creates acess token with refresh token (always inside a transaction)
func grantAccessToken(w rest.ResponseWriter, cnf *config.Config, db *gorm.DB, clientID, userID int) {
	tx := db.Begin()

	refreshToken := RefreshToken{
		RefreshToken: uuid.New(),
		ExpiresAt:    time.Now().Add(time.Duration(cnf.RefreshTokenLifetime) * time.Second),
	}
	if err := tx.Create(&refreshToken).Error; err != nil {
		tx.Rollback()
		api.Error(w, "Error saving refresh token", http.StatusInternalServerError)
		return
	}

	var scopes []string
	db.Model(&Scope{}).Where(&Scope{IsDefault: true}).Pluck("scope", &scopes)

	accessToken := AccessToken{
		AccessToken:    uuid.New(),
		ExpiresAt:      time.Now().Add(time.Duration(cnf.AccessTokenLifetime) * time.Second),
		Scope:          strings.Join(scopes, " "),
		ClientID:       clientID,
		RefreshTokenID: refreshToken.ID,
	}
	if userID > 0 {
		accessToken.UserID = userID
	}
	if err := tx.Create(&accessToken).Error; err != nil {
		tx.Rollback()
		api.Error(w, "Error saving access token", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.RefreshToken,
	})
}
