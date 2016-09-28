package oauth

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/RichardKnop/go-oauth2-server/migrations"
)

var (
	list = []migrations.MigrationStage{
		{
			Name:     "oauth_initial",
			Function: migrate0001,
		},
	}
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	return migrations.Migrate(db, list)
}

func migrate0001(db *gorm.DB, name string) error {
	// Create tables
	if err := db.CreateTable(new(Client)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_clients table: %s", err)
	}
	if err := db.CreateTable(new(Scope)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_scopes table: %s", err)
	}
	if err := db.CreateTable(new(Role)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_roles table: %s", err)
	}
	if err := db.CreateTable(new(User)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_users table: %s", err)
	}
	if err := db.CreateTable(new(RefreshToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_refresh_tokens table: %s", err)
	}
	if err := db.CreateTable(new(AccessToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_access_tokens table: %s", err)
	}
	if err := db.CreateTable(new(AuthorizationCode)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_authorization_codes table: %s", err)
	}
	err := db.Model(new(User)).AddForeignKey(
		"role_id", "oauth_roles(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_users.role_id for oauth_roles(id): %s", err)
	}
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(RefreshToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(AccessToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(AccessToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(AuthorizationCode)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.user_id for oauth_users(id): %s", err)
	}

	return nil
}
