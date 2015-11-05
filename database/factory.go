package database

import (
	"fmt"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/jinzhu/gorm"
	// Drivers
	_ "github.com/lib/pq"
)

// NewDatabase returns a gorm.DB struct, gorm.DB.DB() returns a database handle
// see http://golang.org/pkg/database/sql/#DB
func NewDatabase(cnf *config.Config) (*gorm.DB, error) {
	// Postgres
	if cnf.Database.Type == "postgres" {
		// Connection args
		// see https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
		args := fmt.Sprintf(
			"sslmode=disable host=%s port=%d user=%s password='%s' dbname=%s",
			cnf.Database.Host,
			cnf.Database.Port,
			cnf.Database.User,
			cnf.Database.Password,
			cnf.Database.DatabaseName,
		)
		db, err := gorm.Open(cnf.Database.Type, args)
		if err != nil {
			return &db, err
		}
		return &db, nil
	}

	// Database type not supported
	return &gorm.DB{},
		fmt.Errorf("Database type %s not suppported", cnf.Database.Type)
}
