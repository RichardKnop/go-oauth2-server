package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
)

const (
	accessTokenHint  = "access_token"
	refreshTokenHint = "refresh_token"
)

var (
	errTokenMissing     = errors.New("Token missing")
	errTokenHintInvalid = errors.New("Invalid token hint")
)

func (s *Service) introspectToken(w http.ResponseWriter, r *http.Request, client *Client) {
	// Parse the form so r.Form becomes available
	if err := r.ParseForm(); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := r.Form.Get("token")
	if token == "" {
		response.Error(w, errTokenMissing.Error(), http.StatusBadRequest)
		return
	}

	tokenTypeHint := r.Form.Get("token_type_hint")

	if tokenTypeHint == "" {
		tokenTypeHint = accessTokenHint
	}

	var ir *IntrospectResponse

	switch tokenTypeHint {
	case accessTokenHint:
		var ok bool
		ir, ok = s.introspectAccessToken(w, token)
		if !ok {
			ir, _ = s.introspectRefreshToken(w, token, client)
		}
	case refreshTokenHint:
		var ok bool
		ir, ok = s.introspectRefreshToken(w, token, client)
		if !ok {
			ir, _ = s.introspectAccessToken(w, token)
		}
	default:
		response.Error(w, errTokenHintInvalid.Error(), http.StatusBadRequest)
		return
	}

	if ir == nil {
		ir = &IntrospectResponse{}
	}

	response.WriteJSON(w, ir, 200)
}

func (s *Service) IntrospectResponseAccessToken(at *AccessToken) *IntrospectResponse {
	ir := IntrospectResponse{
		Active:    true,
		Scope:     at.Scope,
		TokenType: TokenType,
		ExpiresAt: int(at.ExpiresAt.Unix()),
	}
	if at.Client != nil {
		ir.ClientID = at.Client.Key
	}
	if at.User != nil {
		ir.Username = at.User.Username
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

func (s *Service) IntrospectResponseRefreshToken(rt *RefreshToken) *IntrospectResponse {
	ir := IntrospectResponse{
		Active:    true,
		Scope:     rt.Scope,
		TokenType: TokenType,
		ExpiresAt: int(rt.ExpiresAt.Unix()),
	}
	if rt.Client != nil {
		ir.ClientID = rt.Client.Key
	}
	if rt.User != nil {
		ir.Username = rt.User.Username
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
