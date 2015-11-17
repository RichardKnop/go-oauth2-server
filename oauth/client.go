package oauth

import (
	"errors"
	"log"

	"github.com/RichardKnop/go-oauth2-server/password"
	"github.com/RichardKnop/go-oauth2-server/util"
)

// AuthClient authenticates client
func (s *Service) AuthClient(clientID, secret string) (*Client, error) {
	// Fetch the client
	client, err := s.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.New("Client not found")
	}

	// Verify the secret
	if password.VerifyPassword(client.Secret, secret) != nil {
		return nil, errors.New("Invalid secret")
	}

	return client, nil
}

// CreateClient saves a new client to database
func (s *Service) CreateClient(clientID, secret, redirectURI string) (*Client, error) {
	secretHash, err := password.HashPassword(secret)
	if err != nil {
		return nil, errors.New("Bcrypt error")
	}
	client := &Client{
		ClientID:    clientID,
		Secret:      string(secretHash),
		RedirectURI: util.StringOrNull(redirectURI),
	}
	if err := s.db.Create(client).Error; err != nil {
		log.Print(err)
		return nil, errors.New("Error saving client to database")
	}
	return client, nil
}

// FindClientByClientID looks up a client by client ID
func (s *Service) FindClientByClientID(clientID string) (*Client, error) {
	// Client IDs are case insensitive
	client := new(Client)
	if s.db.Where("LOWER(client_id) = LOWER(?)", clientID).First(client).RecordNotFound() {
		return nil, errors.New("Client not found")
	}
	return client, nil
}
