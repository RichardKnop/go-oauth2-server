package migrations

import "log"

// RunAll executes all database migrations
func RunAll() error {
	log.Print("Running database migrations")

	if err := migrate0000(); err != nil {
		return err
	}

	if err := migrate0001(); err != nil {
		return err
	}

	return nil
}
