package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	GetOauthService() oauth.ServiceInterface

	// Needed for the newRoutes to be able to register handlers
	authorizeForm(w http.ResponseWriter, r *http.Request)
	authorize(w http.ResponseWriter, r *http.Request)
	loginForm(w http.ResponseWriter, r *http.Request)
	login(w http.ResponseWriter, r *http.Request)
	logout(w http.ResponseWriter, r *http.Request)
	registerForm(w http.ResponseWriter, r *http.Request)
	register(w http.ResponseWriter, r *http.Request)
}
