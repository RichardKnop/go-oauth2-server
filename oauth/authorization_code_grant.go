package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/errors"

	"github.com/ant0ine/go-json-rest/rest"
)

func (s *service) authorizationCodeGrant(w rest.ResponseWriter, r *rest.Request, client *Client) {
	code := r.FormValue("code")

	// Fetch an auth code from the database
	authCode := AuthCode{}
	if s.db.Where(&AuthCode{
		Code:     code,
		ClientID: clientIDOrNull(client),
	}).First(&authCode).RecordNotFound() {
		errors.Error(w, "Auth code not found", http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := s.grantAccessToken(&authCode.Client, &authCode.User, authCode.Scope)
	if err != nil {
		errors.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	writeAccessToken(w, s.cnf.AccessTokenLifetime, accessToken, refreshToken)
}
