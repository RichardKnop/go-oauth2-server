package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/jinzhu/gorm"
)

// Service struct keeps config and db objects to avoid passing them around
type Service struct {
	cnf *config.Config
	db  *gorm.DB
}

var s *Service

// NewService starts a new Service instance
func NewService(cnf *config.Config, db *gorm.DB) *Service {
	s = &Service{cnf: cnf, db: db}
	return s
}

// GetService returns internal Service instance
func GetService() *Service {
	return s
}
