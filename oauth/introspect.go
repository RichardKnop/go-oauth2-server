package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
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

func (s *Service) introspectToken(w http.ResponseWriter, r *http.Request, client *Client) {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Form.Get("token")
	if token == "" {
		response.Error(w, ErrTokenMissing.Error(), http.StatusBadRequest)
		return
	}

	tokenTypeHint := r.Form.Get("token_type_hint")

	if tokenTypeHint == "" {
		tokenTypeHint = AccessTokenHint
	}

	var ir *IntrospectResponse

	switch tokenTypeHint {
	case AccessTokenHint:
		var ok bool
		ir, ok = s.introspectAccessToken(w, token)
		if !ok {
			ir, _ = s.introspectRefreshToken(w, token, client)
		}
	case RefreshTokenHint:
		var ok bool
		ir, ok = s.introspectRefreshToken(w, token, client)
		if !ok {
			ir, _ = s.introspectAccessToken(w, token)
		}
	default:
		response.Error(w, ErrTokenHintInvalid.Error(), http.StatusBadRequest)
		return
	}

	if ir == nil {
		ir = &IntrospectResponse{}
	}

	response.WriteJSON(w, ir, 200)
}

// IntrospectResponseAccessToken ...
func (s *Service) IntrospectResponseAccessToken(at *AccessToken) *IntrospectResponse {
	ir := IntrospectResponse{
		Active:    true,
		Scope:     at.Scope,
		TokenType: TokenType,
		ExpiresAt: int(at.ExpiresAt.Unix()),
	}
	if at.ClientID.Valid {
		c := new(Client)
		notFound := s.db.Select("key").First(c, at.ClientID.Int64).RecordNotFound()
		if !notFound {
			ir.ClientID = c.Key
		}
	}
	if at.UserID.Valid {
		u := new(User)
		notFound := s.db.Select("username").First(u, at.UserID.Int64).RecordNotFound()
		if !notFound {
			ir.Username = u.Username
		}
	}

	return &ir
}

// Introspects give token as access token and returns true if it was successful
func (s *Service) introspectAccessToken(w http.ResponseWriter, token string) (*IntrospectResponse, bool) {
	at, err := s.Authenticate(token)
	if err != nil {
		return nil, false
	}

	return s.IntrospectResponseAccessToken(at), true
}

// IntrospectResponseRefreshToken ...
func (s *Service) IntrospectResponseRefreshToken(rt *RefreshToken) *IntrospectResponse {
	ir := IntrospectResponse{
		Active:    true,
		Scope:     rt.Scope,
		TokenType: TokenType,
		ExpiresAt: int(rt.ExpiresAt.Unix()),
	}
	if rt.ClientID.Valid {
		c := new(Client)
		notFound := s.db.Select("key").First(c, rt.ClientID.Int64).RecordNotFound()
		if !notFound {
			ir.ClientID = c.Key
		}
	}
	if rt.UserID.Valid {
		u := new(User)
		notFound := s.db.Select("username").First(u, rt.UserID.Int64).RecordNotFound()
		if !notFound {
			ir.Username = u.Username
		}
	}
	return &ir
}

// Introspects given token as refresh token and returns true if it was successful
func (s *Service) introspectRefreshToken(w http.ResponseWriter, token string, client *Client) (*IntrospectResponse, bool) {
	rt, err := s.GetValidRefreshToken(token, client)
	if err != nil {
		return nil, false
	}

	return s.IntrospectResponseRefreshToken(rt), true
}
