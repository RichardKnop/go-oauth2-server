package oauth

import (
	"errors"
	"time"
)

var (
	// ErrAccessTokenNotFound ...
	ErrAccessTokenNotFound = errors.New("Access token not found")
	// ErrAccessTokenExpired ...
	ErrAccessTokenExpired = errors.New("Access token expired")
)

// Authenticate checks the access token is valid
func (s *Service) Authenticate(token string) (*AccessToken, error) {
	// Fetch the access token from the database
	accessToken := new(AccessToken)
	notFound := s.db.Where("token = ?", token).
		Preload("Client").Preload("User").First(accessToken).RecordNotFound()

	// Not found
	if notFound {
		return nil, ErrAccessTokenNotFound
	}

	// Check the access token hasn't expired
	if time.Now().After(accessToken.ExpiresAt) {
		return nil, ErrAccessTokenExpired
	}

	return accessToken, nil
}
