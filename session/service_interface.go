package session

// ServiceInterface defines exported methods
type ServiceInterface interface {
	StartSession() error
	GetUserSession() (*UserSession, error)
	SetUserSession(userSession *UserSession) error
	ClearUserSession() error
	SetFlashMessage(msg string) error
	GetFlashMessage() (interface{}, error)
}
