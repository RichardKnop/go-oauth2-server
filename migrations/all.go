package migrations

import "log"

// RunAll executes all database migrations
func RunAll() error {
	log.Print("Running migrations")

	if err := runMigration0000(); err != nil {
		return err
	}

	if err := runMigration0001(); err != nil {
		return err
	}

	return nil
}
