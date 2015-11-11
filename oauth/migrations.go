package oauth

import (
	"fmt"
	"log"

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
	migrationName := "0001_initial"

	migration := new(migrations.Migration)
	found := !db.Where(migrations.Migration{
		Name: migrationName,
	}).First(migration).RecordNotFound()

	if found {
		log.Printf("Skipping %s migration", migrationName)
		return nil
	}

	log.Printf("Running %s migration", migrationName)

	var err error

	// Create clients table
	if err := db.CreateTable(new(Client)).Error; err != nil {
		return fmt.Errorf("Error creating clients table: %s", db.Error)
	}

	// Create scopes table
	if err := db.CreateTable(new(Scope)).Error; err != nil {
		return fmt.Errorf("Error creating scopes table: %s", db.Error)
	}

	// Create users table
	if err := db.CreateTable(new(User)).Error; err != nil {
		return fmt.Errorf("Error creating users table: %s", db.Error)
	}

	// Create refresh_tokens table
	if err := db.CreateTable(new(RefreshToken)).Error; err != nil {
		return fmt.Errorf("Error creating refresh_tokens table: %s", db.Error)
	}

	// Create access_tokens table
	if err := db.CreateTable(new(AccessToken)).Error; err != nil {
		return fmt.Errorf("Error creating access_tokens table: %s", db.Error)
	}

	// Create auth_codes table
	if err := db.CreateTable(new(AuthorizationCode)).Error; err != nil {
		return fmt.Errorf("Error creating authorization_codes table: %s", db.Error)
	}

	// Add foreign key on refresh_tokens.client_id
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"client_id",
		"clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"refresh_tokens.client_id for clients(id): %s", db.Error)
	}

	// Add foreign key on refresh_tokens.user_id
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"user_id",
		"users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"refresh_tokens.user_id for users(id): %s", db.Error)
	}

	// Add foreign key on access_tokens.client_id
	err = db.Model(new(AccessToken)).AddForeignKey(
		"client_id",
		"clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"access_tokens.client_id for clients(id): %s", db.Error)
	}

	// Add foreign key on access_tokens.user_id
	err = db.Model(new(AccessToken)).AddForeignKey(
		"user_id",
		"users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"access_tokens.user_id for users(id): %s", db.Error)
	}

	// Add foreign key on authorization_codes.client_id
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"client_id",
		"clients(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"authorization_codes.client_id for clients(id): %s", db.Error)
	}

	// Add foreign key on authorization_codes.user_id
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"user_id",
		"users(id)",
		"RESTRICT",
		"RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"authorization_codes.user_id for users(id): %s", db.Error)
	}

	// Save a record to migrations table,
	// so we don't rerun this migration again
	migration.Name = migrationName
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
