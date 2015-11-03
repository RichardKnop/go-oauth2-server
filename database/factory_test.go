package database

import (
	"errors"
	"testing"

	"github.com/RichardKnop/go-microservice-example/config"
	"github.com/stretchr/testify/assert"
)

// TODO - this test relies on a database
// func TestNewDatabasePostgres(t *testing.T) {
// 	cnf := config.NewConfig()
// 	_, err := NewDatabase(cnf)
//
// 	assert.Nil(t, err)
// }

func TestNewDatabaseTypeNotSupported(t *testing.T) {
	cnf := &config.Config{
		Database: config.DatabaseConfig{
			Type: "bogus",
		},
	}
	_, err := NewDatabase(cnf)

	if assert.NotNil(t, err) {
		assert.Equal(t, errors.New("Database type bogus not suppported"), err)
	}
}
