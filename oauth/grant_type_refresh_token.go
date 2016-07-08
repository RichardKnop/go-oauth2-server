package oauth

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
	"github.com/RichardKnop/go-oauth2-server/util"
)

var (
	// ErrRequestedScopeCannotBeGreater ...
	ErrRequestedScopeCannotBeGreater = errors.New("Requested scope cannot be greater")
)

func (s *Service) refreshTokenGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Fetch the refresh token
	theRefreshToken, err := s.GetValidRefreshToken(
		r.Form.Get("refresh_token"), // refresh token
		client, // client
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Default to the scope originally granted by the resource owner
	scope := theRefreshToken.Scope

	// If the scope is specified in the request, get the scope string
	if r.Form.Get("scope") != "" {
		scope, err = s.GetScope(r.Form.Get("scope"))
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Requested scope CANNOT include any scope not originally granted
	if !util.SpaceDelimitedStringNotGreater(scope, theRefreshToken.Scope) {
		response.Error(w, ErrRequestedScopeCannotBeGreater.Error(), http.StatusBadRequest)
		return
	}

	// Log in the user
	accessToken, refreshToken, err := s.Login(
		theRefreshToken.Client,
		theRefreshToken.User,
		scope,
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON access token to the response
	accessTokenRespone := &AccessTokenResponse{
		AccessToken:  accessToken.Token,
		ExpiresIn:    s.cnf.Oauth.AccessTokenLifetime,
		TokenType:    TokenType,
		Scope:        accessToken.Scope,
		RefreshToken: refreshToken.Token,
	}
	if accessToken.User != nil {
		accessTokenRespone.UserID = accessToken.User.MetaUserID
	}
	response.WriteJSON(w, accessTokenRespone, 200)
}
