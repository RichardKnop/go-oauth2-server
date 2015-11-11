package accounts

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/jinzhu/gorm"

	"errors"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	oauthService *oauth.Service // oauth service dependency injection
}

var s *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB, oauthService *oauth.Service) *Service {
	s = &Service{cnf: cnf, db: db, oauthService: oauthService}
	return s
}

// GetService returns internal Service instance
func GetService() *Service {
	return s
}

// Register ...
func (s *Service) Register(username, password string) (*oauth.User, error) {
	if s.oauthService.UserExists(username) {
		return nil, errors.New("Username already taken")
	}

	return s.oauthService.CreateUser(username, password)
}

// Login ...
func (s *Service) Login(username, password string) (*oauth.User, error) {
	return s.oauthService.AuthUser(username, password)
}
