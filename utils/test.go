package utils

import (
	"log"

	"github.com/RichardKnop/go-microservice-example/migrations"
	"github.com/jinzhu/gorm"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

var testDBCreated bool
var testDB *gorm.DB

// TestSetUp prepares test database
func TestSetUp() *gorm.DB {
	if testDBCreated {
		return testDB
	}

	testDB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	testDBCreated = true
	migrations.MigrateAll(&testDB)
	return &testDB
}

// TestTearDown empties the test database
func TestTearDown(db *gorm.DB) {
	db.Exec("DELETE FROM SELECT name FROM sqlite_master WHERE type IS 'table'")
}
