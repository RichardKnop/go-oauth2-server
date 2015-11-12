package oauth

import (
	"errors"
	"time"
)

func (s *Service) getOrCreateRefreshToken(client *Client, user *User, scope string) (*RefreshToken, error) {
	// Try to fetch an existing refresh token first
	refreshToken := new(RefreshToken)
	notFound := s.db.Where(RefreshToken{
		ClientID: clientIDOrNull(client),
		UserID:   userIDOrNull(user),
	}).First(refreshToken).RecordNotFound()

	// Check if the token is expired, if found
	var expired bool
	if !notFound {
		expired = time.Now().After(refreshToken.ExpiresAt)
	}

	// If the refresh token has expired, delete it
	if expired {
		s.db.Delete(refreshToken)
	}

	// Create a new refresh token if it expired or was not found
	if expired || notFound {
		refreshToken = newRefreshToken(
			s.cnf.Oauth.RefreshTokenLifetime,
			client,
			user,
			scope,
		)
		if err := s.db.Create(refreshToken).Error; err != nil {
			return nil, errors.New("Error saving refresh token")
		}
	}

	return refreshToken, nil
}
