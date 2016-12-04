package models

import (
	"fmt"

	"github.com/RichardKnop/go-oauth2-server/util/migrations"
	"github.com/jinzhu/gorm"
)

var (
	list = []migrations.MigrationStage{
		{
			Name:     "initial",
			Function: migrate0001,
		},
	}
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	return migrations.Migrate(db, list)
}

func migrate0001(db *gorm.DB, name string) error {
	//-------------
	// OAUTH models
	//-------------

	// Create tables
	if err := db.CreateTable(new(OauthClient)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_clients table: %s", err)
	}
	if err := db.CreateTable(new(OauthScope)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_scopes table: %s", err)
	}
	if err := db.CreateTable(new(OauthRole)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_roles table: %s", err)
	}
	if err := db.CreateTable(new(OauthUser)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_users table: %s", err)
	}
	if err := db.CreateTable(new(OauthRefreshToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_refresh_tokens table: %s", err)
	}
	if err := db.CreateTable(new(OauthAccessToken)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_access_tokens table: %s", err)
	}
	if err := db.CreateTable(new(OauthAuthorizationCode)).Error; err != nil {
		return fmt.Errorf("Error creating oauth_authorization_codes table: %s", err)
	}
	err := db.Model(new(OauthUser)).AddForeignKey(
		"role_id", "oauth_roles(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_users.role_id for oauth_roles(id): %s", err)
	}
	err = db.Model(new(OauthRefreshToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(OauthRefreshToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_refresh_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(OauthAccessToken)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(OauthAccessToken)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_access_tokens.user_id for oauth_users(id): %s", err)
	}
	err = db.Model(new(OauthAuthorizationCode)).AddForeignKey(
		"client_id", "oauth_clients(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.client_id for oauth_clients(id): %s", err)
	}
	err = db.Model(new(OauthAuthorizationCode)).AddForeignKey(
		"user_id", "oauth_users(id)",
		"RESTRICT", "RESTRICT",
	).Error
	if err != nil {
		return fmt.Errorf("Error creating foreign key on "+
			"oauth_authorization_codes.user_id for oauth_users(id): %s", err)
	}

	return nil
}
