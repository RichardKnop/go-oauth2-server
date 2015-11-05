package oauth2

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// TestSuite ...
type TestSuite struct {
	suite.Suite
	DB *gorm.DB
}

// SetupTest ...
func (suite *TestSuite) SetupTest() {
	if suite.DB != nil {
		return
	}
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	suite.DB = &db
	MigrateAll(&db)
}

// TearDown ...
func (suite *TestSuite) TearDown() {
	suite.DB.Exec("DELETE FROM SELECT name FROM sqlite_master WHERE type IS 'table'")
}
