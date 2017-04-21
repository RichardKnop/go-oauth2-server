package redis

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/gorilla/sessions"
	redisStore "gopkg.in/boj/redistore.v1"
)

// ConnectionConfigType is used to connect to the redis db
type ConnectionConfigType struct {
	Size           int
	Network        string
	Address        string
	Password       string
	SessionSecrets [][]byte
}

// SessionOptionsType defines the options for the sessions
type SessionOptionsType struct {
	Path     string
	MaxAge   int
	HTTPOnly bool
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

var (
	// ConnectionConfig ...
	ConnectionConfig ConnectionConfigType
	// SessionOptions ...
	SessionOptions SessionOptionsType
	// DefaultSize ...
	DefaultSize = 10
	// DefaultNetwork ...
	DefaultNetwork = "tcp"
	// DefaultAddress ...
	DefaultAddress = ":6379"
	// DefaultPassword ...
	DefaultPassword = ""
	// DefaultSessionSecrets ...
	DefaultSessionSecrets = "The secret" // cnf.Session.Secret
	// DefaultSessionOptionsPath ...
	DefaultSessionOptionsPath = "/"
	// DefaultSessionOptionsMaxAge ...
	DefaultSessionOptionsMaxAge = 0
	// DefaultSessionOptionsHTTPOnly ...
	DefaultSessionOptionsHTTPOnly = true
	// SessionService the service being exported
	SessionService CustomSessionServiceType
)

func init() {
	ConnectionConfig.Size = DefaultSize
	ConnectionConfig.Network = DefaultNetwork
	ConnectionConfig.Address = DefaultAddress
	ConnectionConfig.Password = DefaultPassword
	ConnectionConfig.SessionSecrets = make([][]byte, 1)
	ConnectionConfig.SessionSecrets[0] = []byte(DefaultSessionSecrets)
	SessionOptions.Path = DefaultSessionOptionsPath
	SessionOptions.MaxAge = DefaultSessionOptionsMaxAge
	SessionOptions.HTTPOnly = DefaultSessionOptionsHTTPOnly
}

// NewPluginService starts the redis connection and sets the session options
func NewPluginService() *CustomSessionServiceType {
	fmt.Println("In InitPluginService")
	store, err := redisStore.NewRediStore(ConnectionConfig.Size, ConnectionConfig.Network, ConnectionConfig.Address, ConnectionConfig.Password, ConnectionConfig.SessionSecrets...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("In InitPluginService after store")

	return &CustomSessionServiceType{
		// Session cookie storage
		sessionStore: store,
		// Session options
		sessionOptions: &sessions.Options{
			Path:     SessionOptions.Path,
			MaxAge:   SessionOptions.MaxAge,
			HttpOnly: SessionOptions.HTTPOnly,
		},
	}
}

func (c *CustomSessionServiceType) GetSessionService() session.ServiceInterface {
	fmt.Println("In GetSessionService", c.sessionOptions.Path)
	return c
}

// Close stops the redis connection
func (c *CustomSessionServiceType) Close() {
	fmt.Println("In Close")
	c.sessionStore.Close()
}

// SetSessionService custom SetSessionStore
func (c *CustomSessionServiceType) SetSessionService(r *http.Request, w http.ResponseWriter) {
	fmt.Println("In SetSessionService")
	c.r = r
	c.w = w
}

// StartSession custom StartSession
func (c *CustomSessionServiceType) StartSession() error {
	fmt.Println("In StartSession", c.sessionOptions.Path)
	// Get a session.
	session, err := c.sessionStore.Get(c.r, session.UserSessionKey)
	if err != nil {
		return err
	}
	fmt.Println("In StartSession after session")

	c.session = session
	fmt.Println("In StartSession after set session")
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
	fmt.Println("b4 flashes in plugins")
	if flashes := c.session.Flashes(); len(flashes) > 0 {
		fmt.Println("Plugins flashes", flashes)
		// We need to save the session, otherwise the flash message won't be removed
		c.session.Save(c.r, c.w)
		return flashes[0], nil
	}

	// No flash messages in the stack
	return nil, nil
}
