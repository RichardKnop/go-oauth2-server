package migrate

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// Migration ...
type Migration struct {
	ID        int
	Name      string `sql:"size:255"`
	CreatedAt time.Time
}

// Bootstrap creates "migrations" table
// to keep track of already run database migrations
func Bootstrap(db *gorm.DB) error {
	migrationName := "0000_bootstrap"

	migration := Migration{}
	if err := db.LogMode(false).Where(&Migration{Name: migrationName}).First(&migration).Error; err != nil {
		log.Printf("Running %s migration", migrationName)

		// Create migrations table
		if err := db.CreateTable(&Migration{}).Error; err != nil {
			return fmt.Errorf("Error creating migrations table: %s", db.Error)
		}

		// Save a record to migrations table,
		// so we don't rerun this migration again
		migration.Name = migrationName
		if err := db.Create(&migration).Error; err != nil {
			return fmt.Errorf("Error saving record to migrations table: %s", err)
		}
	} else {
		log.Printf("Skipping %s migration", migrationName)
	}

	return nil
}
