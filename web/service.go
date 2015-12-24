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
