package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
)

func (s *Service) clientCredentialsGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		client,                          // client
		new(User),                       // empty user
		s.cnf.Oauth.AccessTokenLifetime, // expires in
		scope, // scope
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		client,                           // client
		new(User),                        // empty user
		s.cnf.Oauth.RefreshTokenLifetime, // expires in
		scope, // scope
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON access token to the response
	accessTokenRespone := &AccessTokenResponse{
		ID:           accessToken.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    s.cnf.Oauth.AccessTokenLifetime,
		TokenType:    TokenType,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.Token,
	}
	response.WriteJSON(w, accessTokenRespone, 200)
}
