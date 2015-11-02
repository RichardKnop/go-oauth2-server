package database

import (
	"errors"
	"testing"

	"github.com/RichardKnop/go-microservice-example/config"
)

func TestNewDatabasePostgres(t *testing.T) {
	cnf := config.NewConfig()
	_, err := NewDatabase(cnf)

	if err != nil {
		t.Errorf("err = %s, wanted nil", err.Error())
	}
}

func TestNewDatabaseTypeNotSupported(t *testing.T) {
	cnf := &config.Config{
		Database: config.DatabaseConfig{
			Type: "bogus",
		},
	}
	_, err := NewDatabase(cnf)

	expectedErr := errors.New("Database type bogus not suppported")

	if err == nil {
		t.Errorf("err = nil, wanted %v", expectedErr)
	}

	if err.Error() != expectedErr.Error() {
		t.Errorf("err.Error() = %s, wanted %s", err.Error(), expectedErr.Error())
	}
}
