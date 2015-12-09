package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/response"
	"github.com/RichardKnop/go-oauth2-server/util"
)

func (s *Service) authorizationCodeGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	// Fetch the authorization code
	authorizationCode, err := s.getValidAuthorizationCode(
		r.Form.Get("code"), // authorization code
		client,             // client
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect URI must match if it was used to obtain the authorization code
	if util.StringOrNull(r.Form.Get("redirect_uri")) != authorizationCode.RedirectURI {
		response.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		authorizationCode.Client, // client
		authorizationCode.User,   // user
		authorizationCode.Scope,  // scope
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		authorizationCode.Client, // client
		authorizationCode.User,   // user
		authorizationCode.Scope,  // scope
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the authorization code
	s.db.Unscoped().Delete(&authorizationCode)

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
