package web

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/gorilla/sessions"
)

// sessionService wraps session functionality
type sessionService struct {
	r              *http.Request
	w              http.ResponseWriter
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
}

// userSession has user data stored in a session after logging in
type userSession struct {
	userID       int
	username     string
	accessToken  string
	refreshToken string
}

func init() {
	// Register a new datatype for storage in sessions
	gob.Register(new(userSession))
}

// newSessionService starts a new SessionService instance
func newSessionService(cnf *config.Config, r *http.Request, w http.ResponseWriter) *sessionService {
	return &sessionService{
		r: r,
		w: w,
		// Session cookie storage
		sessionStore: sessions.NewCookieStore([]byte(cnf.Session.Secret)),
		// Session options
		sessionOptions: &sessions.Options{
			Path:     cnf.Session.Path,
			MaxAge:   cnf.Session.MaxAge,
			HttpOnly: cnf.Session.HTTPOnly,
		},
	}
}

// Initialises a new session by name
func (s *sessionService) initSession(name string) error {
	session, err := s.sessionStore.Get(s.r, name)
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

// Sets user session after loggin in or updates it after refreshing tokens
func (s *sessionService) setUser(user *userSession) {
	s.session.Values["user"] = user
	s.session.Save(s.r, s.w)
}

// Gets user session
func (s *sessionService) getUser() (*userSession, error) {
	// Retrieve our user session struct and type-assert it
	var user = new(userSession)
	user, ok := s.session.Values["user"].(*userSession)
	if !ok {
		return nil, errors.New("User session type assertion error")
	}
	return user, nil
}

// Sets a flash message, useful for displaying an error after 302 redirection
func (s *sessionService) setFlashMessage(msg string) {
	s.session.AddFlash(msg)
	s.session.Save(s.r, s.w)
}

// Returns a flash message previously added with setFlashMessage or nil
func (s *sessionService) getFlashMessage() interface{} {
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		s.session.Save(s.r, s.w)
		return flashes[0]
	}
	return nil
}
