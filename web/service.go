package web

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/session"
)

// Service struct keeps variables for reuse
type Service struct {
	cnf            *config.Config
	oauthService   oauth.ServiceInterface
	sessionService session.ServiceInterface
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, oauthService oauth.ServiceInterface, sessionService session.ServiceInterface) *Service {
	return &Service{
		cnf:            cnf,
		oauthService:   oauthService,
		sessionService: sessionService,
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// GetOauthService returns oauth.Service instance
func (s *Service) GetOauthService() oauth.ServiceInterface {
	return s.oauthService
}

// GetSessionService returns session.Service instance
func (s *Service) GetSessionService() session.ServiceInterface {
	return s.sessionService
}

// Close stops any running services
func (s *Service) Close() {}

func (s *Service) setSessionService(r *http.Request, w http.ResponseWriter) {
	s.sessionService.SetSessionService(r, w)
}
