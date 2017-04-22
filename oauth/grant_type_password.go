package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth/tokentypes"
)

var (
	// ErrInvalidUsernameOrPassword ...
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
)

func (s *Service) passwordGrant(r *http.Request, client *models.OauthClient) (*AccessTokenResponse, error) {
	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		return nil, err
	}

	// Authenticate the user
	user, err := s.AuthUser(r.Form.Get("username"), r.Form.Get("password"))
	if err != nil {
		// For security reasons, return a general error message
		return nil, ErrInvalidUsernameOrPassword
	}

	// Log in the user
	accessToken, refreshToken, err := s.Login(client, user, scope)
	if err != nil {
		return nil, err
	}

	// Create response
	accessTokenResponse, err := NewAccessTokenResponse(
		accessToken,
		refreshToken,
		s.cnf.Oauth.AccessTokenLifetime,
		tokentypes.Bearer,
	)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}
