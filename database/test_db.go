package database

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/RichardKnop/go-fixtures"
	"github.com/RichardKnop/go-oauth2-server/migrations"
	"github.com/jinzhu/gorm"
	// sqlite driver
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
func CreateTestDatabasePostgres(dbUser, dbName string, migrationFunctions []func(*gorm.DB) error, fixtureFiles []string) (*gorm.DB, error) {

	// Postgres test database
	db, err := rebuildDatabasePostgres(dbUser, dbName)
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
func rebuildDatabasePostgres(dbUser, dbName string) (*gorm.DB, error) {

	dropPostgresDB(dbUser, dbName)

	if err := createPostgresDB(dbUser, dbName); err != nil {
		return nil, err
	}

	return openPostgresDB(dbUser, dbName)
}

func openPostgresDB(dbUser, dbName string) (*gorm.DB, error) {
	// Init a new postgres test database connection
	db, err := gorm.Open("postgres",
		fmt.Sprintf(
			"sslmode=disable host=localhost port=5432 user=%s password='' dbname=%s",
			dbUser,
			dbName,
		),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createPostgresDB(dbUser, dbName string) error {
	// Create a new test database
	createDbCmd := fmt.Sprintf("createdb -U %s %s", dbUser, dbName)
	log.Println(createDbCmd)
	out, err := exec.Command("sh", "-c", createDbCmd).Output()
	if err != nil {
		log.Printf("%v", string(out))
		return err
	}
	return nil
}

func dropPostgresDB(dbUser, dbName string) {
	// Delete the current database if it exists
	dropDbCmd := fmt.Sprintf("dropdb --if-exists -U %s %s", dbUser, dbName)
	fmt.Println(dropDbCmd)
	exec.Command("sh", "-c", dropDbCmd).Output()
}
