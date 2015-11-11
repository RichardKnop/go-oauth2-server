package oauth

import (
	"errors"
	"time"
)

func (s *Service) grantAccessToken(client *Client, user *User, scope string) (*AccessToken, *RefreshToken, error) {
	// Delete expired access tokens
	s.deleteExpiredAccessTokens(client, user)

	// Create a new access token
	accessToken := newAccessToken(
		s.cnf.Oauth.AccessTokenLifetime,
		client,
		user,
		scope,
	)
	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, nil, errors.New("Error saving access token")
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.getOrCreateRefreshToken(client, user, scope)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) deleteExpiredAccessTokens(client *Client, user *User) {
	s.db.Where(AccessToken{
		ClientID: clientIDOrNull(client),
		UserID:   userIDOrNull(user),
	}).Where("expires_at <= ?", time.Now()).Delete(new(AccessToken))
}

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
