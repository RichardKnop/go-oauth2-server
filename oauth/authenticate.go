package oauth

import (
	"errors"
	"time"
)

var (
	errAccessTokenNotFound = errors.New("Access token not found")
	errAccessTokenExpired  = errors.New("Access token expired")
)

// Authenticate checks the access token is valid
func (s *Service) Authenticate(token string) (*AccessToken, error) {
	// Fetch the access token from the database
	accessToken := new(AccessToken)
	notFound := s.db.Where("token = ?", token).
		Preload("Client").Preload("User").First(accessToken).RecordNotFound()

	// Not found
	if notFound {
		return nil, errAccessTokenNotFound
	}

	// Check the access token hasn't expired
	if time.Now().After(accessToken.ExpiresAt) {
		return nil, errAccessTokenExpired
	}

	return accessToken, nil
}
