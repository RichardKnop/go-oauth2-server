package service

import (
	"net/http"
	"time"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GrantAccessToken ...
func GrantAccessToken(w rest.ResponseWriter, r *rest.Request) {
	// clientID, clientPassword, ok := r.BasicAuth()
	// if !ok {
	// 	w.Header().Set("WWW-Authenticate", "Basic realm=Bearer")
	// 	rest.Error(w, "Unautorized", http.StatusUnauthorized)
	// 	return
	// }

	grantType := r.FormValue("grant_type")

	if grantType == "password" {
		passwordGrant(w, r)
		return
	}

	rest.Error(w, "Invalid grant type", http.StatusBadRequest)
}

func passwordGrant(w rest.ResponseWriter, r *rest.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	cnf := config.NewConfig()

	db, err := database.NewDatabase(cnf)
	if err != nil {
		rest.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}

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

	// TODO proper JSON structure
	w.WriteJson(&accessToken)
}
