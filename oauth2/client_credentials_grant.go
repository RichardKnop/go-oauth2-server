package oauth2

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

func clientCredentialsGrant(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB, client *Client) {
	scope, err := getScope(db, r.FormValue("scope"))
	if err != nil {
		api.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := grantAccessToken(cnf, db, client, nil, scope)
	if err != nil {
		api.Error(w, err.Error(), http.StatusInternalServerError)
	}

	respondWithAccessToken(w, cnf, accessToken)
}
