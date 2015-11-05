package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Handles all OAuth 2.0 grant types
func tokensHandler(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	grantType := r.FormValue("grant_type")

	supportedGrantTypes := map[string]bool{
		"client_credentials": true,
		"password":           true,
		"refresh_token":      true,
	}

	if !supportedGrantTypes[grantType] {
		rest.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	if grantType == "password" {
		password(w, r, cnf, db)
	}

	if grantType == "client_credentials" {
		clientCredentials(w, r, cnf, db)
	}

	if grantType == "refresh_token" {
		refreshToken(w, r, cnf, db)
	}
}

// Grants user credentials access token
func password(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	username, password, ok := r.BasicAuth()
	if !ok {
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

	user := User{}
	if db.Where(&User{Username: username}).First(&user).RecordNotFound() {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	grantAccessToken(w, cnf, db, -1, user.ID)
}

// Grants client credentials access token
func clientCredentials(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		clientID = r.FormValue("client_id")
		clientSecret = r.FormValue("client_secret")
	}

	client := Client{}
	if db.Where(&Client{ClientID: clientID}).First(&client).RecordNotFound() {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(clientSecret)); err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	grantAccessToken(w, cnf, db, client.ID, -1)
}

// Refreshes access token
func refreshToken(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	token := r.FormValue("refresh_token")

	refreshToken := RefreshToken{}
	if db.Where(&RefreshToken{RefreshToken: token}).First(&refreshToken).RecordNotFound() {
		rest.Error(w, "Refresh token not found", http.StatusBadRequest)
		return
	}

	if refreshToken.ExpiresAt.After(time.Now()) {
		rest.Error(w, "Refresh token expired", http.StatusBadRequest)
		return
	}

	accessToken := AccessToken{}
	if db.Where(&AccessToken{RefreshTokenID: refreshToken.ID}).First(&accessToken).RecordNotFound() {
		rest.Error(w, "Access token with refresh token not found", http.StatusBadGateway)
		return
	}

	// delete old access / refresh token?

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
		rest.Error(w, "Error saving refresh token", http.StatusInternalServerError)
		return
	}

	var scopes []Scope
	db.Where(&Scope{IsDefault: true}).Find(&scopes)

	accessToken := AccessToken{
		AccessToken:    uuid.New(),
		ExpiresAt:      time.Now().Add(time.Duration(cnf.AccessTokenLifetime) * time.Second),
		RefreshTokenID: refreshToken.ID,
		Scopes:         scopes,
	}
	if clientID > 0 {
		accessToken.ClientID = clientID
	}
	if userID > 0 {
		accessToken.UserID = userID
	}
	if err := tx.Create(&accessToken).Error; err != nil {
		tx.Rollback()
		rest.Error(w, "Error saving access token", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	scopeStrings := make([]string, len(accessToken.Scopes))
	for _, scope := range accessToken.Scopes {
		scopeStrings = append(scopeStrings, scope.Scope)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         strings.Join(scopeStrings, " "),
		"refresh_token": refreshToken.RefreshToken,
	})
}
