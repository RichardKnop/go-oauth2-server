package database

import (
	"fmt"
	"log"
	"time"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/jinzhu/gorm"

	// Drivers
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
)

func init() {
	gorm.NowFunc = func() time.Time {
		return time.Now().UTC()
	}
}

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
			return db, err
		}

		// Max idle connections
		db.DB().SetMaxIdleConns(cnf.Database.MaxIdleConns)

		// Max open connections
		db.DB().SetMaxOpenConns(cnf.Database.MaxOpenConns)

		// Database logging
		db.LogMode(cnf.IsDevelopment)

		return db, nil
	}
	if cnf.Database.Type == "mysql" {
		var (
			err                                               error
			dbType, dbName, user, password, host, tablePrefix string
			port                                              int
		)
		dbType = cnf.Database.Type
		dbName = cnf.Database.DatabaseName
		user = cnf.Database.User
		password = cnf.Database.Password
		host = cnf.Database.Host
		port = cnf.Database.Port

		db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			user,
			password,
			host,
			port,
			dbName,
		))

		if err != nil {
			log.Println(err)
		}

		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return tablePrefix + defaultTableName
		}

		db.SingularTable(true)
		// Max idle connections
		db.DB().SetMaxIdleConns(cnf.Database.MaxIdleConns)

		// Max open connections
		db.DB().SetMaxOpenConns(cnf.Database.MaxOpenConns)

		// Database logging
		db.LogMode(cnf.IsDevelopment)
		return db, nil
	}
	// Database type not supported
	return nil, fmt.Errorf("Database type %s not suppported", cnf.Database.Type)
}
