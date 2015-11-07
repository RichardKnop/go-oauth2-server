package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func clientCredentialsGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	// Get the scope string
	scope, err := getScope(db, r.FormValue("scope"))
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, err := grantAccessToken(cnf, db, client, nil, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	respondWithAccessToken(w, cnf, accessToken)
}
