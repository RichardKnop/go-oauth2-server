package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/oauth/tokentypes"
)

const (
	// AccessTokenHint ...
	AccessTokenHint = "access_token"
	// RefreshTokenHint ...
	RefreshTokenHint = "refresh_token"
)

var (
	// ErrTokenMissing ...
	ErrTokenMissing = errors.New("Token missing")
	// ErrTokenHintInvalid ...
	ErrTokenHintInvalid = errors.New("Invalid token hint")
)

func (s *Service) introspectToken(r *http.Request, client *Client) (*IntrospectResponse, error) {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	// Get token from the query
	token := r.Form.Get("token")
	if token == "" {
		return nil, ErrTokenMissing
	}

	// Get token type hint from the query
	tokenTypeHint := r.Form.Get("token_type_hint")

	// Default to access token hint
	if tokenTypeHint == "" {
		tokenTypeHint = AccessTokenHint
	}

	switch tokenTypeHint {
	case AccessTokenHint:
		accessToken, err := s.Authenticate(token)
		if err != nil {
			return nil, err
		}
		return s.NewIntrospectResponseFromAccessToken(accessToken)
	case RefreshTokenHint:
		refreshToken, err := s.GetValidRefreshToken(token, client)
		if err != nil {
			return nil, err
		}
		return s.NewIntrospectResponseFromRefreshToken(refreshToken)
	default:
		return nil, ErrTokenHintInvalid
	}
}

// NewIntrospectResponseFromAccessToken ...
func (s *Service) NewIntrospectResponseFromAccessToken(accessToken *AccessToken) (*IntrospectResponse, error) {
	var introspectResponse = &IntrospectResponse{
		Active:    true,
		Scope:     accessToken.Scope,
		TokenType: tokentypes.Bearer,
		ExpiresAt: int(accessToken.ExpiresAt.Unix()),
	}

	if accessToken.ClientID.Valid {
		client := new(Client)
		notFound := s.db.Select("key").First(client, accessToken.ClientID.Int64).
			RecordNotFound()
		if notFound {
			return nil, ErrClientNotFound
		}
		introspectResponse.ClientID = client.Key
	}

	if accessToken.UserID.Valid {
		user := new(User)
		notFound := s.db.Select("username").First(user, accessToken.UserID.Int64).
			RecordNotFound()
		if notFound {
			return nil, ErrUserNotFound
		}
		introspectResponse.Username = user.Username
	}

	return introspectResponse, nil
}

// NewIntrospectResponseFromRefreshToken ...
func (s *Service) NewIntrospectResponseFromRefreshToken(refreshToken *RefreshToken) (*IntrospectResponse, error) {
	var introspectResponse = &IntrospectResponse{
		Active:    true,
		Scope:     refreshToken.Scope,
		TokenType: tokentypes.Bearer,
		ExpiresAt: int(refreshToken.ExpiresAt.Unix()),
	}

	if refreshToken.ClientID.Valid {
		client := new(Client)
		notFound := s.db.Select("key").First(client, refreshToken.ClientID.Int64).
			RecordNotFound()
		if notFound {
			return nil, ErrClientNotFound
		}
		introspectResponse.ClientID = client.Key
	}

	if refreshToken.UserID.Valid {
		user := new(User)
		notFound := s.db.Select("username").First(user, refreshToken.UserID.Int64).
			RecordNotFound()
		if notFound {
			return nil, ErrUserNotFound
		}
		introspectResponse.Username = user.Username
	}

	return introspectResponse, nil
}
