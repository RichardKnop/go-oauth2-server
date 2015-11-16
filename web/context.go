package web

import (
	"errors"
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/session"
	"github.com/gorilla/context"
)

type contextKey int

const sessionServiceKey contextKey = 0
const clientKey contextKey = 1

var (
	errSessionServiceNotPresent = errors.New("Session service not present in the request context.")
	errClientNotPresent         = errors.New("Client not present in the request context.")
)

// Returns *session.Service from the request context
func getSessionService(r *http.Request) (*session.Service, error) {
	val, ok := context.GetOk(r, sessionServiceKey)
	if !ok {
		return nil, errSessionServiceNotPresent
	}

	sessionService, ok := val.(*session.Service)
	if !ok {
		return nil, errSessionServiceNotPresent
	}

	return sessionService, nil
}

// Returns *oauth.Client from the request context
func getClient(r *http.Request) (*oauth.Client, error) {
	val, ok := context.GetOk(r, clientKey)
	if !ok {
		return nil, errClientNotPresent
	}

	client, ok := val.(*oauth.Client)
	if !ok {
		return nil, errClientNotPresent
	}

	return client, nil
}
