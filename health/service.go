package health

import (
	"github.com/jinzhu/gorm"
)

// Service struct keeps db object to avoid passing it around
type Service struct {
	db *gorm.DB
}

// InitHealthService starts a new Service instance
func (s *Service) InitHealthService(db *gorm.DB) {
	s = &Service{db: db}
}
