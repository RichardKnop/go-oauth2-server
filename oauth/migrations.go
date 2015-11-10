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

	migration := migrations.Migration{}
	if err := db.Where(&migrations.Migration{Name: migrationName}).First(&migration).Error; err != nil {
		log.Printf("Running %s migration", migrationName)

		// Create clients table
		if err := db.CreateTable(&Client{}).Error; err != nil {
			return fmt.Errorf("Error creating clients table: %s", db.Error)
		}

		// Create scopes table
		if err := db.CreateTable(&Scope{}).Error; err != nil {
			return fmt.Errorf("Error creating scopes table: %s", db.Error)
		}

		// Create users table
		if err := db.CreateTable(&User{}).Error; err != nil {
			return fmt.Errorf("Error creating users table: %s", db.Error)
		}

		// Create refresh_tokens table
		if err := db.CreateTable(&RefreshToken{}).Error; err != nil {
			return fmt.Errorf("Error creating refresh_tokens table: %s", db.Error)
		}

		// Create access_tokens table
		if err := db.CreateTable(&AccessToken{}).Error; err != nil {
			return fmt.Errorf("Error creating access_tokens table: %s", db.Error)
		}

		// Create auth_codes table
		if err := db.CreateTable(&AuthorizationCode{}).Error; err != nil {
			return fmt.Errorf("Error creating authorization_codes table: %s", db.Error)
		}

		// Add foreign key on refresh_tokens.client_id
		if err := db.Model(&RefreshToken{}).AddForeignKey("client_id", "clients(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on refresh_tokens.client_id for clients(id): %s", db.Error)
		}

		// Add foreign key on refresh_tokens.user_id
		if err := db.Model(&RefreshToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on refresh_tokens.user_id for users(id): %s", db.Error)
		}

		// Add foreign key on access_tokens.client_id
		if err := db.Model(&AccessToken{}).AddForeignKey("client_id", "clients(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on access_tokens.client_id for clients(id): %s", db.Error)
		}

		// Add foreign key on access_tokens.user_id
		if err := db.Model(&AccessToken{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on access_tokens.user_id for users(id): %s", db.Error)
		}

		// Add foreign key on authorization_codes.client_id
		if err := db.Model(&AuthorizationCode{}).AddForeignKey("client_id", "clients(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on authorization_codes.client_id for clients(id): %s", db.Error)
		}

		// Add foreign key on authorization_codes.user_id
		if err := db.Model(&AuthorizationCode{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
			return fmt.Errorf("Error creating foreign key on authorization_codes.user_id for users(id): %s", db.Error)
		}

		// Save a record to migrations table,
		// so we don't rerun this migration again
		migration.Name = migrationName
		if err := db.Create(&migration).Error; err != nil {
			return fmt.Errorf("Error saving record to migrations table: %s", err)
		}
	} else {
		log.Printf("Skipping %s migration", migrationName)
	}

	return nil
}
