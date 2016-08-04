package oauth

import (
	"time"
)

// TokenType is default type of generated tokens.
const TokenType = "Bearer"

// GrantAccessToken deletes old tokens and grants a new access token
func (s *Service) GrantAccessToken(client *Client, user *User, expiresIn int, scope string) (*AccessToken, error) {
	// Delete expired access tokens
	if err := s.DeleteExpiredAccessTokens(client, user); err != nil {
		return nil, err
	}

	// Create a new access token
	accessToken := NewAccessToken(client, user, expiresIn, scope)
	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, err
	}
	accessToken.Client = client
	accessToken.User = user

	return accessToken, nil
}

// DeleteExpiredAccessTokens deletes expired access tokens
func (s *Service) DeleteExpiredAccessTokens(client *Client, user *User) error {
	query := s.db.Unscoped().Where("client_id = ?", client.ID)
	if user != nil && user.ID > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	return query.Where("expires_at <= ?", time.Now()).Delete(new(AccessToken)).Error
}
