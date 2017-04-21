package redis

import (
	"errors"
	"log"
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/gorilla/sessions"
	redisStore "gopkg.in/boj/redistore.v1"
)

// ConfigType is used to connect to the redis db
type ConfigType struct {
	Size           int
	Network        string
	Address        string
	Password       string
	SessionSecrets [][]byte
}

// CustomSessionServiceType extends session.ServiceInterface
type CustomSessionServiceType struct {
	session.ServiceInterface
	sessionStore   *redisStore.RediStore
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
}

// NewService starts the redis connection and sets the session options
func NewService(cnf *config.Config, redisConfig ConfigType) *CustomSessionServiceType {
	store, err := redisStore.NewRediStore(redisConfig.Size, redisConfig.Network, redisConfig.Address, redisConfig.Password, redisConfig.SessionSecrets...)
	if err != nil {
		log.Fatal(err)
	}

	return &CustomSessionServiceType{
		// Session cookie storage
		sessionStore: store,
		// Session options
		sessionOptions: &sessions.Options{
			Path:     cnf.Session.Path,
			MaxAge:   cnf.Session.MaxAge,
			HttpOnly: cnf.Session.HTTPOnly,
		},
	}
}

// Close stops the redis connection
func (c *CustomSessionServiceType) Close() {
	c.sessionStore.Close()
}

// SetSessionService custom SetSessionStore
func (c *CustomSessionServiceType) SetSessionService(r *http.Request, w http.ResponseWriter) {
	c.r = r
	c.w = w
}

// StartSession custom StartSession
func (c *CustomSessionServiceType) StartSession() error {
	// Get a session.
	session, err := c.sessionStore.Get(c.r, session.UserSessionKey)
	if err != nil {
		return err
	}

	c.session = session
	return nil
}

// GetUserSession custom GetUserSession
func (c *CustomSessionServiceType) GetUserSession() (*session.UserSession, error) {
	// Make sure StartSession has been called
	if c.session == nil {
		return nil, session.ErrSessonNotStarted
	}

	// Retrieve our user session struct and type-assert it
	userSession, ok := c.session.Values[session.UserSessionKey].(*session.UserSession)
	if !ok {
		return nil, errors.New("User session type assertion error")
	}

	return userSession, nil
}

// SetUserSession custom SetUserSession
func (c *CustomSessionServiceType) SetUserSession(userSession *session.UserSession) error {
	// Make sure StartSession has been called
	if c.session == nil {
		return session.ErrSessonNotStarted
	}

	// Set a new user session
	c.session.Values[session.UserSessionKey] = userSession
	return c.session.Save(c.r, c.w)
}

// ClearUserSession custom ClearUserSession
func (c *CustomSessionServiceType) ClearUserSession() error {
	c.session.Options.MaxAge = -1
	return c.session.Save(c.r, c.w)
}

// SetFlashMessage custom SetFlashMessage
func (c *CustomSessionServiceType) SetFlashMessage(msg string) error {
	// Make sure StartSession has been called
	if c.session == nil {
		return session.ErrSessonNotStarted
	}

	// Add the flash message
	c.session.AddFlash(msg)
	return c.session.Save(c.r, c.w)
}

// GetFlashMessage custom GetFlashMessage
func (c *CustomSessionServiceType) GetFlashMessage() (interface{}, error) {
	// Make sure StartSession has been called
	if c.session == nil {
		return nil, session.ErrSessonNotStarted
	}

	// Get the last flash message from the stack
	if flashes := c.session.Flashes(); len(flashes) > 0 {
		// We need to save the session, otherwise the flash message won't be removed
		c.session.Save(c.r, c.w)
		return flashes[0], nil
	}

	// No flash messages in the stack
	return nil, nil
}
