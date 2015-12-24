package health

import (
	"github.com/jinzhu/gorm"
)

// Service struct keeps db object to avoid passing it around
type Service struct {
	db *gorm.DB
}

// NewService starts a new Service instance
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}
