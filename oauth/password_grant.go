package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func passwordGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	requestedScope := r.FormValue("scope")

	// Authenticate the user
	user, err := authUser(r.Request, db)
	if err != nil {
		api.UnauthorizedError(w, err.Error())
		return
	}

	// Get the scope string
	scope, err := getScope(db, requestedScope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a new access token
	accessToken, refreshToken, err := grantAccessToken(cnf, db, client, user, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Write the access token to a JSON response
	respondWithAccessToken(w, cnf, accessToken, refreshToken)
}
