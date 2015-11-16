package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
)

func logout(w http.ResponseWriter, r *http.Request) {
	// Initialise the session service
	sessionService := session.NewService(
		theService.cnf,
		r,
		w,
		theService.oauthService,
	)
	if err := sessionService.InitUserSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionService.LogOut()
	http.Redirect(w, r, "/web/login", http.StatusFound)
}
