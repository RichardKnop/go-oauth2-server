package oauth

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Checks client credentials from basic HTTP authentication
func authClient(r *http.Request, db *gorm.DB) (*Client, error) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("Client authentication required")
	}

	client := Client{}
	// Client IDs are case insensitive
	if db.Where("LOWER(client_id) = LOWER(?)", clientID).First(&client).RecordNotFound() {
		return nil, errors.New("Client authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(clientSecret)); err != nil {
		return nil, errors.New("Client authentication failed")
	}

	return &client, nil
}

// Checks user credentials from posted form data
func authUser(r *http.Request, db *gorm.DB) (*User, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user := User{}
	// Usernames are case insensitive
	if db.Where("LOWER(username) = LOWER(?)", username).First(&user).RecordNotFound() {
		return nil, errors.New("User authentication failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("User authentication failed")
	}

	return &user, nil
}
