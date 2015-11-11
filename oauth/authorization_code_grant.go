package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/json"
)

func (s *Service) authorizationCodeGrant(w http.ResponseWriter, r *http.Request, client *Client) {
	code := r.FormValue("code")

	// Fetch an auth code from the database
	authorizationCode := AuthorizationCode{}
	if s.db.Where(&AuthorizationCode{
		Code:     code,
		ClientID: clientIDOrNull(client),
	}).Preload("Client").Preload("User").First(&authorizationCode).RecordNotFound() {
		json.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(&authorizationCode.Client, &authorizationCode.User, authorizationCode.Scope)
	if err != nil {
		json.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Delete the authorization code
	s.db.Delete(&authorizationCode)

	// Write the access token to a JSON response
	writeJSON(w, s.cnf.Oauth.AccessTokenLifetime, accessToken, refreshToken)
}
