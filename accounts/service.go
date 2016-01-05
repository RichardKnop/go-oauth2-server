package accounts

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	oauthService oauth.ServiceInterface // oauth service dependency injection
}

var s *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, oauthService oauth.ServiceInterface) *Service {
	return &Service{
		cnf:          cnf,
		db:           db,
		oauthService: oauthService,
	}
}

// GetOauthService returns oauth.Service instance
func (s *Service) GetOauthService() oauth.ServiceInterface {
	return s.oauthService
}
