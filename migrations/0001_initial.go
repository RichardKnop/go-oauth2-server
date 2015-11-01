package migrations

import (
	"fmt"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
	"github.com/RichardKnop/go-microservice-example/service"
)

func runMigration0001() error {
	migrationName := "0001_initial"

	// Config factory
	cnf := config.Factory()

	// Database connection factory
	db, err := database.Factory(cnf)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err)
	}

	var migration *Migration
	if db.Where(&Migration{Name: migrationName}).First(&migration).RecordNotFound() {
		// Create users table
		if err := db.CreateTable(&service.User{}).Error; err != nil {
			return fmt.Errorf("Error creating users table: %s", db.Error)
		}

		// Create refresh_tokens table
		if err := db.CreateTable(&service.RefreshToken{}).Error; err != nil {
			return fmt.Errorf("Error creating refresh_tokens table: %s", db.Error)
		}

		// Create access_tokens table
		if err := db.CreateTable(&service.AccessToken{}).Error; err != nil {
			return fmt.Errorf("Error creating access_tokens table: %s", db.Error)
		}

		// Save a record to migrations table,
		// so we don't rerun this migration again
		migration.Name = migrationName
		if err := db.Create(&migration).Error; err != nil {
			return fmt.Errorf("Error saving record to migrations table: %s", err)
		}
	}

	return nil
}
