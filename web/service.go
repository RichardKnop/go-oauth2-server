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

var theService *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, oauthService *oauth.Service) *Service {
	theService = &Service{
		cnf:          cnf,
		oauthService: oauthService,
	}
	return theService
}

// GetService returns internal Service instance
func GetService() *Service {
	return theService
}
