package commands

import (
	"io/ioutil"

	"github.com/RichardKnop/go-fixtures"
)

// LoadData loads fixtures
func LoadData(paths []string) error {
	cnf, db, err := initConfigDB(true, false)
	defer db.Close()
	if err != nil {
		return err
	}

	// Iterate over fixtures paths
	for _, path := range paths {
		// Read the contents of the fixture
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Load the data from the fixture into the database
		if err := fixtures.Load(data, db.DB(), cnf.Database.Type); err != nil {
			return err
		}
	}

	return nil
}
