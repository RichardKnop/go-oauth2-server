package web

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
)

// Service struct keeps variables for reuse
type Service struct {
	cnf          *config.Config
	oauthService oauth.ServiceInterface
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, oauthService oauth.ServiceInterface) *Service {
	return &Service{
		cnf:          cnf,
		oauthService: oauthService,
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
