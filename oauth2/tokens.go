package oauth2

import (
	"net/http"
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
	client, err := authClient(r.Request, db)
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
	user, err := authUser(r.Request, db)
	if err != nil {
		api.UnauthorizedError(w, err.Error())
		return
	}

	grantAccessToken(w, cnf, db, client, user, r.FormValue("scope"))
}

// Grants client credentials access token
func clientCredentials(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	grantAccessToken(w, cnf, db, client, nil, r.FormValue("scope"))
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
	if db.Model(&accessToken).Related(&refreshToken).RecordNotFound() {
		api.Error(w, "Access token not found", http.StatusBadGateway)
		return
	}

	requestedScope := r.FormValue("scope")
	// Requested scope CANNOT include any scope not originally granted
	if !scopeNotGreater(requestedScope, accessToken.Scope) {
		api.Error(w, "Invalid scope", http.StatusBadGateway)
		return
	}

	// Delete old access / refresh token
	db.Delete(&refreshToken)
	db.Delete(&accessToken)

	grantAccessToken(w, cnf, db, &accessToken.Client, &accessToken.User, requestedScope)
}

// Creates acess token with refresh token (always inside a transaction)
func grantAccessToken(w rest.ResponseWriter, cnf *config.Config, db *gorm.DB, client *Client, user *User, requestedScope string) {
	scope, err := getScope(db, requestedScope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken := AccessToken{
		AccessToken: uuid.New(),
		ExpiresAt:   time.Now().Add(time.Duration(cnf.AccessTokenLifetime) * time.Second),
		Scope:       scope,
		Client:      *client,
		RefreshToken: RefreshToken{
			RefreshToken: uuid.New(),
			ExpiresAt:    time.Now().Add(time.Duration(cnf.RefreshTokenLifetime) * time.Second),
		},
	}
	if user != nil {
		accessToken.User = *user
	}
	if err := db.Create(&accessToken).Error; err != nil {
		api.Error(w, "Error saving access token", http.StatusInternalServerError)
		return
	}

	// TODO - should we delete old access tokens for this client / user?

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": accessToken.RefreshToken.RefreshToken,
	})
}
