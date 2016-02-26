package database

import (
	"os"

	"github.com/RichardKnop/go-fixtures"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/jinzhu/gorm"
)

// rebuildDatabase attempts to delete an existing in memory
// database and rebuild it, returning a pointer to it
func rebuildDatabase(dbPath string) (*gorm.DB, error) {
	// Delete the current database if it exists
	os.Remove(dbPath)

	// Init in-memory test database
	inMemoryDB, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &inMemoryDB, nil
}

// CreateTestDatabase recreates the test database and
// runs migrations and fixtures as passed in, returning
// a pointer to the database
func CreateTestDatabase(dbPath string, migrationFunctions []func(*gorm.DB) error, fixtureFiles []string) (*gorm.DB, error) {

	// Init in-memory test database
	inMemoryDB, err := rebuildDatabase(dbPath)
	if err != nil {
		return nil, err
	}

	// Run all migrations
	migrations.MigrateAll(inMemoryDB, migrationFunctions)

	if err = fixtures.LoadFiles(fixtureFiles, inMemoryDB.DB(), "sqlite"); err != nil {
		return nil, err
	}

	return inMemoryDB, nil
}
