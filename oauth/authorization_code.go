package oauth

import (
	"errors"
	"time"
)

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
