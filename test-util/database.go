package testutil

import (
	"fmt"
	"os"

	"github.com/RichardKnop/go-fixtures"
	"github.com/RichardKnop/go-oauth2-server/util/migrations"
	"github.com/jinzhu/gorm"

	// Drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

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

	// Load data from data
	if err = fixtures.LoadFiles(fixtureFiles, inMemoryDB.DB(), "sqlite"); err != nil {
		return nil, err
	}

	return inMemoryDB, nil
}

// CreateTestDatabasePostgres is similar to CreateTestDatabase but it uses
// Postgres instead of sqlite, this is needed for testing packages that rely
// on some Postgres specifuc features (such as table inheritance)
func CreateTestDatabasePostgres(dbHost, dbUser, dbName string, migrationFunctions []func(*gorm.DB) error, fixtureFiles []string) (*gorm.DB, error) {

	// Postgres test database
	db, err := rebuildDatabasePostgres(dbHost, dbUser, dbName)
	if err != nil {
		return nil, err
	}

	// Run all migrations
	migrations.MigrateAll(db, migrationFunctions)

	// Load data from data
	if err = fixtures.LoadFiles(fixtureFiles, db.DB(), "postgres"); err != nil {
		return nil, err
	}

	return db, nil
}

// rebuildDatabase attempts to delete an existing in memory
// database and rebuild it, returning a pointer to it
func rebuildDatabase(dbPath string) (*gorm.DB, error) {
	// Delete the current database if it exists
	os.Remove(dbPath)

	// Init a new in-memory test database connection
	inMemoryDB, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return inMemoryDB, nil
}

// rebuildDatabase attempts to delete an existing Postgres
// database and rebuild it, returning a pointer to it
func rebuildDatabasePostgres(dbHost, dbUser, dbName string) (*gorm.DB, error) {
	db, err := openPostgresDB(dbHost, dbUser, "template1")
	if err != nil {
		return nil, err
	}

	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName)).Error; err != nil {
		return nil, err
	}

	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error; err != nil {
		return nil, err
	}

	return openPostgresDB(dbHost, dbUser, dbName)
}

func openPostgresDB(dbHost, dbUser, dbName string) (*gorm.DB, error) {
	// Init a new postgres test database connection
	db, err := gorm.Open("postgres",
		fmt.Sprintf(
			"sslmode=disable host=%s port=5432 user=%s password='' dbname=%s",
			dbHost,
			dbUser,
			dbName,
		),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
