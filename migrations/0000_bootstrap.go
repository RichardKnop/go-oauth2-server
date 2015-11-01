package migrations

import (
	"fmt"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/RichardKnop/go-microservice-example/database"
)

func runMigration0000() error {
	// Config factory
	cnf := config.Factory()

	// Database connection factory
	db, err := database.Factory(cnf)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %s", err)
	}

	// Create migrations table
	if err := db.CreateTable(&Migration{}).Error; err != nil {
		return fmt.Errorf("Error creating migrations table: %s", db.Error)
	}

	return nil
}
