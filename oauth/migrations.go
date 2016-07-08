package oauth

import (
	"fmt"

	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/jinzhu/gorm"
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	if err := migrate0001(db); err != nil {
		return err
	}

	return nil
}

// Migrate0001 creates OAuth 2.0 schema
func migrate0001(db *gorm.DB) error {
	migrationName := "oauth_initial"

	migration := new(migrations.Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
		return nil
	}

	logger.Infof("Running %s migration", migrationName)

	var err error

	// Create oauth_clients table
	if err := db.CreateTable(new(Client)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_clients table: %s", err)
	}

	// Create oauth_scopes table
	if err := db.CreateTable(new(Scope)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_scopes table: %s", err)
	}

	// Create oauth_users table
	if err := db.CreateTable(new(User)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_users table: %s", err)
	}

	// Create oauth_refresh_tokens table
	if err := db.CreateTable(new(RefreshToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_refresh_tokens table: %s", err)
	}

	// Create oauth_access_tokens table
	if err := db.CreateTable(new(AccessToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_access_tokens table: %s", err)
	}

	// Create oauth_authorization_codes table
	if err := db.CreateTable(new(AuthorizationCode)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_authorization_codes table: %s", err)
	}

	// Add foreign key on oauth_refresh_tokens.client_id
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"client_id",
		"oauth_clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.client_id for oauth_clients(id): %s", err)
	}

	// Add foreign key on oauth_refresh_tokens.user_id
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"user_id",
		"oauth_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.user_id for oauth_users(id): %s", err)
	}

	// Add foreign key on oauth_access_tokens.client_id
	err = db.Model(new(AccessToken)).AddForeignKey(
		"client_id",
		"oauth_clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.client_id for oauth_clients(id): %s", err)
	}

	// Add foreign key on oauth_access_tokens.user_id
	err = db.Model(new(AccessToken)).AddForeignKey(
		"user_id",
		"oauth_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.user_id for oauth_users(id): %s", err)
	}

	// Add foreign key on oauth_authorization_codes.client_id
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"client_id",
		"oauth_clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.client_id for oauth_clients(id): %s", err)
	}

	// Add foreign key on oauth_authorization_codes.user_id
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"user_id",
		"oauth_users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.user_id for oauth_users(id): %s", err)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
