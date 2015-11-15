package oauth

import (
	"errors"
	"log"
)

// AuthClient authenticates client
func (s *Service) AuthClient(clientID, secret string) (*Client, error) {
	// Fetch the client
	client, err := s.FindClientByClientID(clientID)
	if err != nil {
		return nil, errors.New("Client not found")
	}

	// Verify the secret
	if verifyPassword(client.Secret, secret) != nil {
		return nil, errors.New("Invalid secret")
	}

	return client, nil
}

// CreateClient saves a new client to database
func (s *Service) CreateClient(clientID, secret string) (*Client, error) {
	secretHash, err := hashPassword(secret)
	if err != nil {
		return nil, errors.New("Bcrypt error")
	}
	client := &Client{
		ClientID: clientID,
		Secret:   string(secretHash),
	}
	if err := s.db.Create(client).Error; err != nil {
		log.Print(err)
		return nil, errors.New("Error saving client to database")
	}
	return client, nil
}

// FindClientByClientID looks up a client by client id
func (s *Service) FindClientByClientID(clientID string) (*Client, error) {
	// Client IDs are case insensitive
	client := new(Client)
	if s.db.Where("LOWER(client_id) = LOWER(?)", clientID).First(client).RecordNotFound() {
		return nil, errors.New("Client not found")
	}
	return client, nil
}
