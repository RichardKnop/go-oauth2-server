package oauth

import (
	"errors"
	"time"
)

// GrantAccessToken grants a new access token
func (s *Service) GrantAccessToken(client *Client, user *User, scope string) (*AccessToken, *RefreshToken, error) {
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
