package oauth

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

// TokenType is default type of generated tokens.
const TokenType = "Bearer"

// GrantAccessToken deletes old tokens and grants a new access token
func (s *Service) GrantAccessToken(client *Client, user *User, expiresIn int, scope string) (*AccessToken, error) {
	// Delete expired access tokens
	s.DeleteExpiredAccessTokens(client, user)

	// Create a new access token
	accessToken := NewAccessToken(client, user, expiresIn, scope)
	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// DeleteExpiredAccessTokens deletes expired access tokens
func (s *Service) DeleteExpiredAccessTokens(client *Client, user *User) {
	clientID := util.PositiveIntOrNull(int64(client.ID))
	userID := util.PositiveIntOrNull(0) // user ID can be NULL
	if user != nil {
		userID = util.PositiveIntOrNull(int64(user.ID))
	}
	s.db.Unscoped().Where(AccessToken{
		ClientID: clientID,
		UserID:   userID,
	}).Where("expires_at <= ?", time.Now()).Delete(new(AccessToken))
}
