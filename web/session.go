package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/gorilla/sessions"
)

// A name used to identify the user session
var sessionName = "user_session"

// sessionService wraps session functionality
type sessionService struct {
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
}

// newSessionService starts a new SessionService instance
func newSessionService(cnf *config.Config) *sessionService {
	return &sessionService{
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

// Initialises session named
func (s *sessionService) initSession(r *http.Request) error {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

// Sets a flash message, useful for displaying an error after 302 redirection
func (s *sessionService) setFlashMessage(msg string, r *http.Request, w http.ResponseWriter) {
	s.session.AddFlash(msg)
	s.session.Save(r, w)
}

// Returns a flash message previously added with setFlashMessage or nil
func (s *sessionService) getFlashMessage(r *http.Request, w http.ResponseWriter) interface{} {
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		s.session.Save(r, w)
		return flashes[0]
	}
	return nil
}
