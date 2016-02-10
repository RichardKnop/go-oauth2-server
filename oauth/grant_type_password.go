package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
)

var (
	errUserAuthenticationRequired = errors.New("User authentication required")
)

func (s *Service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Get user credentials from form data
	username := r.Form.Get("username") // usually an email
	password := r.Form.Get("password")

	// Authenticate the user
	user, err := s.AuthUser(username, password)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, errUserAuthenticationRequired.Error())
		return
	}

	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		client, // client
		user,   // user
		scope,  // scope
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		client, // client
		user,   // user
		scope,  // scope
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
