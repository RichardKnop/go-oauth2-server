package commands

import (
	"github.com/RichardKnop/go-fixtures"
)

// LoadData loads fixtures
func LoadData(paths []string) error {
	cnf, db, err := initConfigDB(true, false)
	if err != nil {
		return err
	}
	defer db.Close()
	return fixtures.LoadFiles(paths, db.DB(), cnf.Database.Type)
}
