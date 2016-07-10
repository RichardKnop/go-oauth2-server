package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

var (
	// ErrAuthorizationCodeNotFound ...
	ErrAuthorizationCodeNotFound = errors.New("Authorization code not found")
	// ErrAuthorizationCodeExpired ...
	ErrAuthorizationCodeExpired = errors.New("Authorization code expired")
)

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *Client, user *User, expiresIn int, redirectURI, scope string) (*AuthorizationCode, error) {
	// Create a new authorization code
	authorizationCode := NewAuthorizationCode(client, user, expiresIn, redirectURI, scope)
	if err := s.db.Create(authorizationCode).Error; err != nil {
		return nil, err
	}

	return authorizationCode, nil
}

// GetValidAuthorizationCode returns a valid non expired authorization code
func (s *Service) GetValidAuthorizationCode(code string, client *Client) (*AuthorizationCode, error) {
	// Fetch the auth code from the database
	authorizationCode := new(AuthorizationCode)
	clientID := util.PositiveIntOrNull(int64(client.ID))
	notFound := s.db.Where(AuthorizationCode{ClientID: clientID}).
		Where("code = ?", code).Preload("Client").Preload("User").
		First(authorizationCode).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAuthorizationCodeNotFound
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, ErrAuthorizationCodeExpired
	}

	return authorizationCode, nil
}
