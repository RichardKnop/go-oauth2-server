package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/gorilla/context"
)

// parseFormMiddleware parses the form so r.Form becomes available
type parseFormMiddleware struct{}

// ServeHTTP as per the negroni.Handler interface
func (m *parseFormMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	next(w, r)
}

// guestMiddleware just initialises session
type guestMiddleware struct {
	service *Service
}

// newGuestMiddleware creates a new guestMiddleware instance
func newGuestMiddleware(service *Service) *guestMiddleware {
	return &guestMiddleware{service: service}
}

// ServeHTTP as per the negroni.Handler interface
func (m *guestMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Initialise the session service
	sessionService := session.NewService(m.service.cnf, r, w)

	// Attempt to start the session
	if err := sessionService.StartSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Set(r, sessionServiceKey, sessionService)

	next(w, r)
}

// loggedInMiddleware initialises session and makes sure the user is logged in
type loggedInMiddleware struct {
	service *Service
}

// newLoggedInMiddleware creates a new loggedInMiddleware instance
func newLoggedInMiddleware(service *Service) *loggedInMiddleware {
	return &loggedInMiddleware{service: service}
}

// ServeHTTP as per the negroni.Handler interface
func (m *loggedInMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Initialise the session service
	sessionService := session.NewService(m.service.cnf, r, w)

	// Attempt to start the session
	if err := sessionService.StartSession(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	context.Set(r, sessionServiceKey, sessionService)

	// Try to get a user session
	userSession, err := sessionService.GetUserSession()
	if err != nil {
		redirectWithQueryString("/web/login", r.URL.Query(), w, r)
		return
	}

	// Authenticate
	if err := m.authenticate(userSession); err != nil {
		redirectWithQueryString("/web/login", r.URL.Query(), w, r)
		return
	}

	// Update the user session
	sessionService.SetUserSession(userSession)

	next(w, r)
}

func (m *loggedInMiddleware) authenticate(userSession *session.UserSession) error {
	// Try to authenticate with the stored access token
	_, err := m.service.oauthService.Authenticate(userSession.AccessToken)
	if err == nil {
		// Access token valid, return
		return nil
	}
	// Access token might be expired, let's try refreshing...

	// Fetch the client
	client, err := m.service.oauthService.FindClientByClientID(
		userSession.ClientID, // client ID
	)
	if err != nil {
		return err
	}

	// Validate the refresh token
	theRefreshToken, err := m.service.oauthService.GetValidRefreshToken(
		userSession.RefreshToken, // refresh token
		client, // client
	)
	if err != nil {
		return err
	}

	// Create a new access token
	accessToken, err := m.service.oauthService.GrantAccessToken(
		theRefreshToken.Client, // client
		theRefreshToken.User,   // user
		theRefreshToken.Scope,  // scope
	)
	if err != nil {
		return err
	}

	// Create or retrieve a refresh token
	refreshToken, err := m.service.oauthService.GetOrCreateRefreshToken(
		theRefreshToken.Client, // client
		theRefreshToken.User,   // user
		theRefreshToken.Scope,  // scope
	)
	if err != nil {
		return err
	}

	userSession.AccessToken = accessToken.Token
	userSession.RefreshToken = refreshToken.Token

	return nil
}

// clientMiddleware takes client_id param from the query string and
// makes a database lookup for a client with the same client ID
type clientMiddleware struct {
	service *Service
}

// newClientMiddleware creates a new clientMiddleware instance
func newClientMiddleware(service *Service) *clientMiddleware {
	return &clientMiddleware{service: service}
}

// ServeHTTP as per the negroni.Handler interface
func (m *clientMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Fetch the client
	client, err := m.service.oauthService.FindClientByClientID(
		r.Form.Get("client_id"), // client ID
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	context.Set(r, clientKey, client)

	next(w, r)
}
