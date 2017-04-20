package plugins

import (
	"log"
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/health"
	"github.com/adam-hanna/go-oauth2-server/oauth"
	"github.com/adam-hanna/go-oauth2-server/session"
	"github.com/adam-hanna/go-oauth2-server/web"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	redisStore "gopkg.in/boj/redistore.v1"
)

// CustomHealthService extends health.ServiceInterface
type CustomHealthService struct {
	health.ServiceInterface
	db *gorm.DB
}

// NewHealthService defines a custom health service if the developer so chooses to implement one
func NewHealthService(db *gorm.DB) *CustomHealthService {
	// YOUR CODE, HERE
	return nil
}

// CustomAuthService extends oauth.ServiceInterface
type CustomAuthService struct {
	oauth.ServiceInterface
	cnf          *config.Config
	db           *gorm.DB
	allowedRoles []string
}

// NewOauthService defines a custom auth service if the developer so chooses to implement one
func NewOauthService(cnf *config.Config, db *gorm.DB) *CustomAuthService {
	// YOUR CODE, HERE
	return nil
}

// CustomSessionService extends session.ServiceInterface
type CustomSessionService struct {
	session.ServiceInterface
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
}

// NewSessionService defines a custom session service if the developer so chooses to implement one
func NewSessionService(cnf *config.Config) *CustomSessionService {
	store, err := redisStore.NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
	if err != nil {
		log.Fatal(err)
	}
	// note @adam-hanna: how to handle this?
	// defer store.Close()

	return &CustomSessionService{
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

// CustomWebService extends web.ServiceInterface
type CustomWebService struct {
	web.ServiceInterface
	cnf            *config.Config
	oauthService   oauth.ServiceInterface
	sessionService session.ServiceInterface
}

// NewWebService defines a custom web service if the developer so chooses to implement one
func NewWebService(cnf *config.Config, oauthService oauth.ServiceInterface, sessionService session.ServiceInterface) *CustomWebService {
	// YOUR CODE, HERE
	return nil
}
