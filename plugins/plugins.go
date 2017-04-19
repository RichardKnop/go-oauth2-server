package plugins

import (
	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/health"
	"github.com/adam-hanna/go-oauth2-server/oauth"
	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/adam-hanna/go-oauth2-server/web"
	"github.com/jinzhu/gorm"
)

// CustomHealthService extends health.ServiceInterface
type CustomHealthService struct {
	health.ServiceInterface
}

// NewHealthService defines a custom health service if the developer so chooses to implement one
func NewHealthService(db *gorm.DB) *CustomHealthService {
	// YOUR CODE, HERE
	return nil
}

// CustomAuthService extends health.ServiceInterface
type CustomAuthService struct {
	oauth.ServiceInterface
}

// NewOauthService defines a custom auth service if the developer so chooses to implement one
func NewOauthService(cnf *config.Config, db *gorm.DB) *CustomAuthService {
	// YOUR CODE, HERE
	return nil
}

// CustomSessionService extends health.ServiceInterface
type CustomSessionService struct {
	session.ServiceInterface
}

// NewSessionService defines a custom session service if the developer so chooses to implement one
func NewSessionService(cnf *config.Config) *CustomSessionService {
	// YOUR CODE, HERE
	return nil
}

// CustomWebService extends health.ServiceInterface
type CustomWebService struct {
	web.ServiceInterface
}

// NewWebService defines a custom web service if the developer so chooses to implement one
func NewWebService(cnf *config.Config, oauthService oauth.ServiceInterface, sessionService session.ServiceInterface) *CustomWebService {
	// YOUR CODE, HERE
	return nil
}
