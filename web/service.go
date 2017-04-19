package web

import (
	"fmt"
	"net/http"

	"github.com/adam-hanna/go-oauth2-server/config"
	"github.com/adam-hanna/go-oauth2-server/oauth"
	"github.com/adam-hanna/go-oauth2-server/session"
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
	fmt.Println("in get config", s.cnf)
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

func (s *Service) setSessionService(r *http.Request, w http.ResponseWriter) {
	fmt.Println("oAuth then session", s.oauthService, s.sessionService)
	s.sessionService.SetSessionService(r, w)
}
