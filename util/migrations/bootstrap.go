package migrations

import (
	"fmt"

	"github.com/RichardKnop/go-oauth2-server/log"
	"github.com/jinzhu/gorm"
)

// Bootstrap creates "migrations" table
// to keep track of already run database migrations
func Bootstrap(db *gorm.DB) error {
	migrationName := "bootstrap_migrations"

	migration := new(Migration)
	// Using Error instead of RecordNotFound because we want to check
	// if the migrations table exists. This is different from later migrations
	// where we query the already create migrations table.
	exists := nil == db.Where("name = ?", migrationName).First(migration).Error

	if exists {
		log.INFO.Printf("Skipping %s migration", migrationName)
		return nil
	}

	log.INFO.Printf("Running %s migration", migrationName)

	// Create migrations table
	if err := db.CreateTable(new(Migration)).Error; err != nil {
		return fmt.Errorf("Error creating migrations table: %s", db.Error)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
