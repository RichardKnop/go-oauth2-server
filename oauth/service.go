package oauth

import (
	"github.com/jinzhu/gorm"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/oauth/roles"
)

// Service struct keeps objects to avoid passing them around
type Service struct {
	cnf          *config.Config
	db           *gorm.DB
	allowedRoles []string
}

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB) *Service {
	return &Service{
		cnf:          cnf,
		db:           db,
		allowedRoles: []string{roles.Superuser, roles.User},
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// RestrictToRoles restricts this service to only specified roles
func (s *Service) RestrictToRoles(allowedRoles ...string) {
	s.allowedRoles = allowedRoles
}

// IsRoleAllowed returns true if the role is allowed to use this service
func (s *Service) IsRoleAllowed(role string) bool {
	for _, allowedRole := range s.allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
