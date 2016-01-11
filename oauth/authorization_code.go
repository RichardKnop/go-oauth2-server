package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

var (
	errAuthorizationCodeNotFound = errors.New("Authorization code not found")
	errAuthorizationCodeExpired  = errors.New("Authorization code expired")
)

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *Client, user *User, redirectURI, scope string) (*AuthorizationCode, error) {
	// Create a new authorization code
	authorizationCode := newAuthorizationCode(
		s.cnf.Oauth.AuthCodeLifetime, // expires in
		client,      // client
		user,        // user
		redirectURI, // redirect URI
		scope,       // scope
	)
	if err := s.db.Create(authorizationCode).Error; err != nil {
		return nil, err
	}

	return authorizationCode, nil
}

// getValidAuthorizationCode returns a valid non expired authorization code
func (s *Service) getValidAuthorizationCode(code string, client *Client) (*AuthorizationCode, error) {
	// Fetch the auth code from the database
	authorizationCode := new(AuthorizationCode)
	notFound := s.db.Where(AuthorizationCode{
		ClientID: util.IntOrNull(int64(client.ID)),
	}).Where("code = ?", code).Preload("Client").Preload("User").
		First(authorizationCode).RecordNotFound()

	// Not found
	if notFound {
		return nil, errAuthorizationCodeNotFound
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, errAuthorizationCodeExpired
	}

	return authorizationCode, nil
}
