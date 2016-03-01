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

	// Log in the user
	accessToken, refreshToken, err := s.Login(client, new(User), scope)
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
