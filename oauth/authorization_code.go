package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
)

var (
	// ErrAuthorizationCodeNotFound ...
	ErrAuthorizationCodeNotFound = errors.New("Authorization code not found")
	// ErrAuthorizationCodeExpired ...
	ErrAuthorizationCodeExpired = errors.New("Authorization code expired")
)

// GrantAuthorizationCode grants a new authorization code
func (s *Service) GrantAuthorizationCode(client *models.OauthClient, user *models.OauthUser, expiresIn int, redirectURI, scope string) (*models.OauthAuthorizationCode, error) {
	// Create a new authorization code
	authorizationCode := models.NewOauthAuthorizationCode(client, user, expiresIn, redirectURI, scope)
	if err := s.db.Create(authorizationCode).Error; err != nil {
		return nil, err
	}
	authorizationCode.Client = client
	authorizationCode.User = user

	return authorizationCode, nil
}

// getValidAuthorizationCode returns a valid non expired authorization code
func (s *Service) getValidAuthorizationCode(code, redirectURI string, client *models.OauthClient) (*models.OauthAuthorizationCode, error) {
	// Fetch the auth code from the database
	authorizationCode := new(models.OauthAuthorizationCode)
	notFound := models.OauthAuthorizationCodePreload(s.db).Where("client_id = ?", client.ID).
		Where("code = ?", code).First(authorizationCode).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAuthorizationCodeNotFound
	}

	// Redirect URI must match if it was used to obtain the authorization code
	if redirectURI != authorizationCode.RedirectURI.String {
		return nil, ErrInvalidRedirectURI
	}

	// Check the authorization code hasn't expired
	if time.Now().After(authorizationCode.ExpiresAt) {
		return nil, ErrAuthorizationCodeExpired
	}

	return authorizationCode, nil
}
