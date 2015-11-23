package session

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/gorilla/sessions"
)

// Service wraps session functionality
type Service struct {
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
}

// UserSession has user data stored in a session after logging in
type UserSession struct {
	ClientID     string
	Username     string
	AccessToken  string
	RefreshToken string
}

var (
	storageSessionName  = "session"
	userSessionKey      = "user"
	errSessonNotStarted = errors.New("Session not started")
)

func init() {
	// Register a new datatype for storage in sessions
	gob.Register(new(UserSession))
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, r *http.Request, w http.ResponseWriter) *Service {
	return &Service{
		// Session cookie storage
		sessionStore: sessions.NewCookieStore([]byte(cnf.Session.Secret)),
		// Session options
		sessionOptions: &sessions.Options{
			Path:     cnf.Session.Path,
			MaxAge:   cnf.Session.MaxAge,
			HttpOnly: cnf.Session.HTTPOnly,
		},
		r: r,
		w: w,
	}
}

// StartSession starts a new session. This method must be called before other
// public methods of this struct as it sets the internal session object
func (s *Service) StartSession() error {
	session, err := s.sessionStore.Get(s.r, storageSessionName)
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

// GetUserSession returns the user session
func (s *Service) GetUserSession() (*UserSession, error) {
	// Make sure StartSession has been called
	if s.session == nil {
		return nil, errSessonNotStarted
	}

	// Retrieve our user session struct and type-assert it
	userSession, ok := s.session.Values[userSessionKey].(*UserSession)
	if !ok {
		return nil, errors.New("User session type assertion error")
	}

	return userSession, nil
}

// SetUserSession saves the user session
func (s *Service) SetUserSession(userSession *UserSession) error {
	// Make sure StartSession has been called
	if s.session == nil {
		return errSessonNotStarted
	}

	// Set a new user session
	s.session.Values[userSessionKey] = userSession
	return s.session.Save(s.r, s.w)
}

// ClearUserSession deletes the user session
func (s *Service) ClearUserSession() error {
	// Make sure StartSession has been called
	if s.session == nil {
		return errSessonNotStarted
	}

	// Delete the user session
	delete(s.session.Values, userSessionKey)
	return s.session.Save(s.r, s.w)
}

// SetFlashMessage sets a flash message,
// useful for displaying an error after 302 redirection
func (s *Service) SetFlashMessage(msg string) error {
	// Make sure StartSession has been called
	if s.session == nil {
		return errSessonNotStarted
	}

	// Add the flash message
	s.session.AddFlash(msg)
	return s.session.Save(s.r, s.w)
}

// GetFlashMessage returns the first flash message
func (s *Service) GetFlashMessage() (interface{}, error) {
	// Make sure StartSession has been called
	if s.session == nil {
		return nil, errSessonNotStarted
	}

	// Get the last flash message from the stack
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		s.session.Save(s.r, s.w)
		return flashes[0], nil
	}

	// No flash messages in the stack
	return nil, nil
}
