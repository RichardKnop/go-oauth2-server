package oauth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// AuthClient authenticates client
func (s *Service) AuthClient(clientID, clientSecret string) (*Client, error) {
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

// AuthUser authenticates user
func (s *Service) AuthUser(username, password string) (*User, error) {
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
