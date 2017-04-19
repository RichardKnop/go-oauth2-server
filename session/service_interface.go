package session

import (
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/config"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	InitService(cnf *config.Config, r *http.Request, w http.ResponseWriter)
	StartSession() error
	GetUserSession() (*UserSession, error)
	SetUserSession(userSession *UserSession) error
	ClearUserSession() error
	SetFlashMessage(msg string) error
	GetFlashMessage() (interface{}, error)
}
