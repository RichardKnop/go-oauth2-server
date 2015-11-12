package web

import (
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/gorilla/sessions"
)

var oauthService *oauth.Service
var sessionStore sessions.Store

// Bootstrap sets internal variables
func Bootstrap(sessionSecret string, s *oauth.Service) {
	oauthService = s
	sessionStore = sessions.NewCookieStore([]byte(sessionSecret))
}
