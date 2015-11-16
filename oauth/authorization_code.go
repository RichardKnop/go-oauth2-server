package oauth

import (
	"errors"
	"time"
)

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *Client, user *User, scope string) (*AuthorizationCode, error) {
	// Create a new authorization code
	authorizationCode := newAuthorizationCode(
		s.cnf.Oauth.AuthCodeLifetime,
		client,
		user,
		scope,
	)
	if err := s.db.Create(authorizationCode).Error; err != nil {
		return nil, errors.New("Error saving authorization code")
	}

	return authorizationCode, nil
}

func (s *Service) getValidAuthorizationCode(code string, client *Client) (*AuthorizationCode, error) {
	// Fetch the auth code from the database
	authorizationCode := new(AuthorizationCode)
	if s.db.Where(AuthorizationCode{
		Code:     code,
		ClientID: clientIDOrNull(client),
	}).Preload("Client").Preload("User").First(authorizationCode).RecordNotFound() {
		return nil, errors.New("Authorization code not found")
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, errors.New("Authorization code expired")
	}

	return authorizationCode, nil
}
