package oauth

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

// TokenType is default type of generated tokens.
const TokenType = "Bearer"

// GrantAccessToken deletes old tokens and grants a new access token
func (s *Service) GrantAccessToken(client *Client, user *User, scope string) (*AccessToken, error) {
	// Delete expired access tokens
	s.deleteExpiredAccessTokens(client, user)

	// Create a new access token
	accessToken := newAccessToken(
		s.cnf.Oauth.AccessTokenLifetime, // expires in
		client, // client
		user,   // user
		scope,  // scope
	)
	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// deleteExpiredAccessTokens deletes expired access tokens
func (s *Service) deleteExpiredAccessTokens(client *Client, user *User) {
	s.db.Unscoped().Where(AccessToken{
		ClientID: util.PositiveIntOrNull(int64(client.ID)),
		UserID:   util.PositiveIntOrNull(int64(user.ID)),
	}).Where("expires_at <= ?", time.Now()).Delete(new(AccessToken))
}
