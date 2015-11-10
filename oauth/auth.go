package oauth

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *service) authClient(r *http.Request) (*Client, error) {
	// Get credentials from basic auth
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("Client credentials required")
	}

	// Fetch the client
	client := Client{}
	// Client IDs are case insensitive
	if s.db.Where("LOWER(client_id) = LOWER(?)", clientID).First(&client).RecordNotFound() {
		return nil, errors.New("Client authentication failed")
	}

	// Check the secret
	if err := bcrypt.CompareHashAndPassword([]byte(client.Secret), []byte(clientSecret)); err != nil {
		return nil, errors.New("Client authentication failed")
	}

	return &client, nil
}

func (s *service) authUser(r *http.Request) (*User, error) {
	// Get credentials from from the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Fetch the user
	user := User{}
	// Usernames are case insensitive
	if s.db.Where("LOWER(username) = LOWER(?)", username).First(&user).RecordNotFound() {
		return nil, errors.New("User authentication failed")
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("User authentication failed")
	}

	return &user, nil
}
