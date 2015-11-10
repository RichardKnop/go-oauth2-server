package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/errors"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

// service struct keeps config and db objects to avoid passing them around
type service struct {
	cnf *config.Config
	db  *gorm.DB
}

// POST /oauth2/api/v1/tokens (handles all OAuth 2.0 grant types)
func (s *service) handleTokens(w rest.ResponseWriter, r *rest.Request) {
	// Check the grant type
	grantTypes := map[string]bool{
		"authorization_code": true,
		"password":           true,
		"client_credentials": true,
		"refresh_token":      true,
	}
	if !grantTypes[r.FormValue("grant_type")] {
		errors.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	// Authenticate the client
	client, err := s.authClient(r.Request)
	if err != nil {
		errors.UnauthorizedError(w, err.Error())
		return
	}

	grants := map[string]func(){
		"authorization_code": func() { s.authorizationCodeGrant(w, r, client) },
		"password":           func() { s.passwordGrant(w, r, client) },
		"client_credentials": func() { s.clientCredentialsGrant(w, r, client) },
		"refresh_token":      func() { s.refreshTokenGrant(w, r, client) },
	}
	grants[r.FormValue("grant_type")]()
}
