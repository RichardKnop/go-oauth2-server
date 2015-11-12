package web

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/gorilla/sessions"
)

var (
	oauthService   *oauth.Service
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
)

// Bootstrap sets internal variables
func Bootstrap(cnf *config.Config, s *oauth.Service) {
	oauthService = s

	// Session options
	sessionOptions = &sessions.Options{
		Path:     cnf.Session.Path,
		MaxAge:   cnf.Session.MaxAge,
		HttpOnly: cnf.Session.HTTPOnly,
	}

	// Session cookie storage
	sessionStore = sessions.NewCookieStore([]byte(cnf.Session.Secret))
}
