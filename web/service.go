package web

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
)

// Service struct keeps variables for reuse
type Service struct {
	cnf          *config.Config
	oauthService *oauth.Service
}

var s *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, oauthService *oauth.Service) *Service {
	s = &Service{
		cnf:          cnf,
		oauthService: oauthService,
	}
	return s
}

// GetService returns internal Service instance
func GetService() *Service {
	return s
}
