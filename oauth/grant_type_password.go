package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
)

var (
	// ErrInvalidUsernameOrPassword ...
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
)

func (s *Service) passwordGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Get user credentials from form data
	username := r.Form.Get("username") // usually an email
	password := r.Form.Get("password")

	// Authenticate the user
	user, err := s.AuthUser(username, password)
	if err != nil {
		// For security reasons, return a general error message
		response.UnauthorizedError(w, ErrInvalidUsernameOrPassword.Error())
		return
	}

	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log in the user
	accessToken, refreshToken, err := s.Login(client, user, scope)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON access token to the response
	accessTokenRespone := &AccessTokenResponse{
		UserID:       user.MetaUserID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    s.cnf.Oauth.AccessTokenLifetime,
		TokenType:    TokenType,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.Token,
	}
	response.WriteJSON(w, accessTokenRespone, 200)
}
