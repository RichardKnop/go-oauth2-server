package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/models"
)

// Login creates an access token and refresh token for a user (logs him/her in)
func (s *Service) Login(client *models.OauthClient, user *models.OauthUser, scope string) (*models.OauthAccessToken, *models.OauthRefreshToken, error) {
	// Return error if user's role is not allowed to use this service
	if !s.IsRoleAllowed(user.RoleID.String) {
		// For security reasons, return a general error message
		return nil, nil, ErrInvalidUsernameOrPassword
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		client,
		user,
		s.cnf.Oauth.AccessTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		client,
		user,
		s.cnf.Oauth.RefreshTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
