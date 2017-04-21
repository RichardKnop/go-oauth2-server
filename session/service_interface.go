package session

import "net/http"

// ServiceInterface defines exported methods
type ServiceInterface interface {
	SetSessionService(r *http.Request, w http.ResponseWriter)
	StartSession() error
	GetUserSession() (*UserSession, error)
	SetUserSession(userSession *UserSession) error
	ClearUserSession() error
	SetFlashMessage(msg string) error
	GetFlashMessage() (interface{}, error)
	Close()
}
