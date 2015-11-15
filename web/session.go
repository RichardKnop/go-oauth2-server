package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/gorilla/sessions"
)

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

func (s *sessionService) initSession(r *http.Request) error {
	session, err := s.sessionStore.Get(r, "areatech")
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

func (s *sessionService) addFlashMessage(msg string, r *http.Request, w http.ResponseWriter) {
	s.session.AddFlash(msg)
	s.session.Save(r, w)
}

func (s *sessionService) getLastFlashMessage(r *http.Request, w http.ResponseWriter) interface{} {
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		s.session.Save(r, w)
		return flashes[0]
	}
	return nil
}
