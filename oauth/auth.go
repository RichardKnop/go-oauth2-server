package oauth

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func authClient(r *http.Request, db *gorm.DB) (*Client, error) {
	// Get credentials from basic auth
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("Client credentials required")
	}

	// Fetch the client
	client := Client{}
	// Client IDs are case insensitive
	if db.Where("LOWER(client_id) = LOWER(?)", clientID).First(&client).RecordNotFound() {
		return nil, errors.New("Client authentication failed")
	}

	// Check the secret
	if err := bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret)); err != nil {
		return nil, errors.New("Client authentication failed")
	}

	return &client, nil
}

func authUser(r *http.Request, db *gorm.DB) (*User, error) {
	// Get credentials from from the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Fetch the user
	user := User{}
	// Usernames are case insensitive
	if db.Where("LOWER(username) = LOWER(?)", username).First(&user).RecordNotFound() {
		return nil, errors.New("User authentication failed")
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("User authentication failed")
	}

	return &user, nil
}
