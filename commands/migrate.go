package commands

import (
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/RichardKnop/go-oauth2-server/oauth"
)

// Migrate runs database migrations
func Migrate() error {
	_, db, err := initConfigDB(true, false)
	defer db.Close()
	if err != nil {
		return err
	}

	// Bootstrap migrations
	if err := migrations.Bootstrap(db); err != nil {
		return err
	}

	// Run migrations for the oauth service
	if err := oauth.MigrateAll(db); err != nil {
		return err
	}

	return nil
}
