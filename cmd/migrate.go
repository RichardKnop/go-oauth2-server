package cmd

import (
	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/util/migrations"
)

// Migrate runs database migrations
func Migrate(configBackend string) error {
	_, db, err := initConfigDB(true, false, configBackend)
	if err != nil {
		return err
	}
	defer db.Close()

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := models.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
