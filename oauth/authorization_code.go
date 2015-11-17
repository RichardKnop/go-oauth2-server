package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *Client, user *User, redirectURI, scope string) (*AuthorizationCode, error) {
	// Create a new authorization code
	authorizationCode := newAuthorizationCode(
		s.cnf.Oauth.AuthCodeLifetime,
		client,
		user,
		redirectURI,
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
		ClientID: util.IntOrNull(client.ID),
	}).Preload("Client").Preload("User").First(authorizationCode).RecordNotFound() {
		return nil, errors.New("Authorization code not found")
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, errors.New("Authorization code expired")
	}

	return authorizationCode, nil
}
