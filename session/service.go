package session

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/gorilla/sessions"
)

// Service wraps session functionality
type Service struct {
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
	oauthService   *oauth.Service
}

// UserSession has user data stored in a session after logging in
type UserSession struct {
	Client       *oauth.Client
	User         *oauth.User
	AccessToken  *oauth.AccessToken
	RefreshToken *oauth.RefreshToken
}

func init() {
	// Register a new datatype for storage in sessions
	gob.Register(new(UserSession))
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, r *http.Request, w http.ResponseWriter, oauthService *oauth.Service) *Service {
	return &Service{
		// Session cookie storage
		sessionStore: sessions.NewCookieStore([]byte(cnf.Session.Secret)),
		// Session options
		sessionOptions: &sessions.Options{
			Path:     cnf.Session.Path,
			MaxAge:   cnf.Session.MaxAge,
			HttpOnly: cnf.Session.HTTPOnly,
		},
		r:            r,
		w:            w,
		oauthService: oauthService,
	}
}

// InitUserSession initialises a new user session
func (s *Service) InitUserSession() error {
	session, err := s.sessionStore.Get(s.r, "user_session")
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

// IsLoggedIn retrieves the user session and authenticates against the oauth
func (s *Service) IsLoggedIn() error {
	// Retrieve our user session struct and type-assert it
	var userSession = new(UserSession)
	userSession, ok := s.session.Values["user"].(*UserSession)
	if !ok {
		return errors.New("User session type assertion error")
	}

	// Try to authenticate with the stored access token
	if err := s.oauthService.Authenticate(userSession.AccessToken.Token); err != nil {
		// Now let's try refreshing the access token with the stored refresh token
		if err := s.refreshToken(userSession); err != nil {
			return err
		}
	}

	return nil
}

// LogIn logs the user in and stores the user session in a cookie
func (s *Service) LogIn(userSession *UserSession) error {
	// Make sure we store any sensitive information in the session
	userSession.Client.Secret = ""
	userSession.User.Password = ""
	// Save the user session
	s.session.Values["user"] = userSession
	return s.session.Save(s.r, s.w)
}

// LogOut clears the user session
func (s *Service) LogOut() {
	delete(s.session.Values, "user")
}

// SetFlashMessage sets a flash message,
// useful for displaying an error after 302 redirection
func (s *Service) SetFlashMessage(msg string) {
	s.session.AddFlash(msg)
	s.session.Save(s.r, s.w)
}

// GetFlashMessage returns the first flash message
func (s *Service) GetFlashMessage() interface{} {
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		s.session.Save(s.r, s.w)
		return flashes[0]
	}
	return nil
}

// Tries to use the stored refresh token to obtain a new access token
func (s *Service) refreshToken(userSession *UserSession) error {
	// Validate the refresh token
	theRefreshToken, err := s.oauthService.ValidateRefreshToken(
		userSession.RefreshToken.Token,
		userSession.Client,
	)
	if err != nil {
		return err
	}

	// Grant an access token
	accessToken, err := s.oauthService.GrantAccessToken(
		userSession.Client,
		userSession.User,
		theRefreshToken.Scope,
	)
	if err != nil {
		return err
	}

	// Get a refresh token
	refreshToken, err := s.oauthService.GetOrCreateRefreshToken(
		userSession.Client,
		userSession.User,
		theRefreshToken.Scope,
	)
	if err != nil {
		return err
	}

	// Update the user session
	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken
	s.session.Values["user"] = userSession
	return s.session.Save(s.r, s.w)
}
