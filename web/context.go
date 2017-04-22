package web

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/gorilla/context"
)

type contextKey int

const (
	sessionServiceKey contextKey = 0
	clientKey         contextKey = 1
)

var (
	// ErrSessionServiceNotPresent ...
	ErrSessionServiceNotPresent = errors.New("Session service not present in the request context")
	// ErrClientNotPresent ...
	ErrClientNotPresent = errors.New("Client not present in the request context")
)

// Returns *session.Service from the request context
func getSessionService(r *http.Request) (session.ServiceInterface, error) {
	val, ok := context.GetOk(r, sessionServiceKey)
	if !ok {
		return nil, ErrSessionServiceNotPresent
	}

	sessionService, ok := val.(session.ServiceInterface)
	if !ok {
		return nil, ErrSessionServiceNotPresent
	}

	return sessionService, nil
}

// Returns *oauth.Client from the request context
func getClient(r *http.Request) (*models.OauthClient, error) {
	val, ok := context.GetOk(r, clientKey)
	if !ok {
		return nil, ErrClientNotPresent
	}

	client, ok := val.(*models.OauthClient)
	if !ok {
		return nil, ErrClientNotPresent
	}

	return client, nil
}
