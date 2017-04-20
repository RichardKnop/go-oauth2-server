package plugins

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/session"
)

// SetSessionService custom SetSessionStore
func (c *CustomSessionService) SetSessionService(r *http.Request, w http.ResponseWriter) {
	c.r = r
	c.w = w
}

// StartSession custom StartSession
func (c *CustomSessionService) StartSession() error {
	// Get a session.
	session, err := c.sessionStore.Get(c.r, session.UserSessionKey)
	if err != nil {
		return err
	}

	c.session = session
	return nil
}

// GetUserSession custom GetUserSession
func (c *CustomSessionService) GetUserSession() (*session.UserSession, error) {
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
func (c *CustomSessionService) SetUserSession(userSession *session.UserSession) error {
	// Make sure StartSession has been called
	if c.session == nil {
		return session.ErrSessonNotStarted
	}

	// Set a new user session
	c.session.Values[session.UserSessionKey] = userSession
	return c.session.Save(c.r, c.w)
}

// ClearUserSession custom ClearUserSession
func (c *CustomSessionService) ClearUserSession() error {
	c.session.Options.MaxAge = -1
	return c.session.Save(c.r, c.w)
}

// SetFlashMessage custom SetFlashMessage
func (c *CustomSessionService) SetFlashMessage(msg string) error {
	// Make sure StartSession has been called
	if c.session == nil {
		return session.ErrSessonNotStarted
	}

	// Add the flash message
	c.session.AddFlash(msg)
	return c.session.Save(c.r, c.w)
}

// GetFlashMessage custom GetFlashMessage
func (c *CustomSessionService) GetFlashMessage() (interface{}, error) {
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
