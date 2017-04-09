package cmd

import (
	"github.com/RichardKnop/go-fixtures"
)

// LoadData loads fixtures
func LoadData(paths []string, configBackend string) error {
	cnf, db, err := initConfigDB(true, false, configBackend)
	if err != nil {
		return err
	}
	defer db.Close()
	return fixtures.LoadFiles(paths, db.DB(), cnf.Database.Type)
}
