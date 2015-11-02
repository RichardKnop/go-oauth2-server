package migrations

import (
	"fmt"
	"log"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
)

// The very first migration creates "migrations" table
// to keep track of already run database migrations
func migrate0000() error {
	migrationName := "0000_bootstrap"

	// Config factory
	cnf := config.NewConfig()

	// Database connection factory
	db, err := database.NewDatabase(cnf)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err)
	}

	migration := &Migration{}
	if err := db.LogMode(false).Where(&Migration{Name: migrationName}).First(migration).Error; err != nil {
		log.Printf("Running %s migration", migrationName)

		// Create migrations table
		if err := db.CreateTable(&Migration{}).Error; err != nil {
			return fmt.Errorf("Error creating migrations table: %s", db.Error)
		}

		// Save a record to migrations table,
		// so we don't rerun this migration again
		migration.Name = migrationName
		if err := db.Create(migration).Error; err != nil {
			return fmt.Errorf("Error saving record to migrations table: %s", err)
		}
	} else {
		log.Printf("Skipping %s migration", migrationName)
	}

	return nil
}
