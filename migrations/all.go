package migrations

import (
	"log"

	"github.com/jinzhu/gorm"
)

// MigrateAll executes all database migrations
func MigrateAll(db *gorm.DB) error {
	log.Print("Running database migrations")

	if err := migrate0000(db); err != nil {
		return err
	}

	if err := migrate0001(db); err != nil {
		return err
	}

	return nil
}
