package oauth

import (
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/jinzhu/gorm"
)

// service struct keeps config and db objects to avoid passing them around
type service struct {
	cnf *config.Config
	db  *gorm.DB
}

var s *service

// InitService starts a new service instance
func InitService(cnf *config.Config, db *gorm.DB) {
	s = &service{cnf: cnf, db: db}
}
