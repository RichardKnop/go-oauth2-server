package oauth

import (
	"errors"
)

// AuthClient authenticates client
func (s *Service) AuthClient(clientID, secret string) (*Client, error) {
	// Fetch the client
	client, err := s.findClientByClientID(clientID)
	if err != nil {
		return nil, errors.New("Client authentication failed")
	}

	// Verify the secret
	if verifyPassword(client.Secret, secret) != nil {
		return nil, errors.New("Client authentication failed")
	}

	return client, nil
}

// CreateClient saves a new client to database
func (s *Service) CreateClient(clientID, secret string) (*Client, error) {
	secretHash, err := hashPassword(secret)
	if err != nil {
		return nil, errors.New("Bcrypt error")
	}
	client := Client{
		ClientID: clientID,
		Secret:   string(secretHash),
	}
	if err := s.db.Create(&client).Error; err != nil {
		return nil, errors.New("Error saving client to database")
	}
	return &client, nil
}

func (s *Service) findClientByClientID(clientID string) (*Client, error) {
	// Client IDs are case insensitive
	client := Client{}
	if s.db.Where("LOWER(client_id) = LOWER(?)", clientID).First(&client).RecordNotFound() {
		return nil, errors.New("Client not found")
	}
	return &client, nil
}
