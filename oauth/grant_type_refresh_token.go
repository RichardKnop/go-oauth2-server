package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
	"github.com/RichardKnop/go-oauth2-server/util"
)

func (s *Service) refreshTokenGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Validate the refresh token
	theRefreshToken, err := s.ValidateRefreshToken(
		r.Form.Get("refresh_token"), // refresh token
		client, // client
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the scope string
	scope, err := s.GetScope(r.Form.Get("scope"))
	if err != nil {
		json.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Requested scope CANNOT include any scope not originally granted
	if !util.SpaceDelimitedStringNotGreater(scope, theRefreshToken.Scope) {
		json.Error(w, "Requested scope cannot be greater", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		theRefreshToken.Client, // client
		theRefreshToken.User,   // user
		scope,                  // scope
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		theRefreshToken.Client, // client
		theRefreshToken.User,   // user
		scope,                  // scope
	)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
