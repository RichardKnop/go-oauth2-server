package migrations

import (
	"fmt"

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
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

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
