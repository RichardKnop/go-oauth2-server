package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokensHandler ...
func TokensHandler(w rest.ResponseWriter, r *rest.Request) {
	clientID, clientPassword, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	grantType := r.FormValue("grant_type")

	supportedGrantTypes := map[string]bool{
		"password": true,
	}

	if !supportedGrantTypes[grantType] {
		rest.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	cnf := config.NewConfig()

	db, err := database.NewDatabase(cnf)
	if err != nil {
		rest.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

	client := Client{}
	if db.Where("client_id = ?", clientID).First(&client).RecordNotFound() {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(clientPassword)); err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if grantType == "password" {
		passwordGrant(w, r, cnf, db)
		return
	}
}

func passwordGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := User{}
	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
		rest.Error(w, "Unautorized", http.StatusUnauthorized)
		return
	}

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

	accessToken := AccessToken{
		AccessToken:    uuid.New(),
		ExpiresAt:      time.Now().Add(time.Duration(cnf.AccessTokenLifetime) * time.Second),
		UserID:         user.ID,
		RefreshTokenID: refreshToken.ID,
	}
	if err := tx.Create(&accessToken).Error; err != nil {
		tx.Rollback()
		rest.Error(w, "Error saving access token", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	var scopes []string
	db.Model(&Scope{}).Where("is_default = ?", true).Pluck("scope", &scopes)

	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         strings.Join(scopes, " "),
		"refresh_token": refreshToken.RefreshToken,
	})
}
