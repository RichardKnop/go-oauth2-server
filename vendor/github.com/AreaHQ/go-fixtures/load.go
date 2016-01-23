package fixtures

import (
	"database/sql"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Load processes a YAML fixture and inserts/updates the database accordingly
func Load(data []byte, db *sql.DB, driver string) error {
	// Unmarshal the YAML data into a []Row slice
	var rows []Row
	if err := yaml.Unmarshal(data, &rows); err != nil {
		return err
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Iterate over rows define in the fixture
	for _, row := range rows {
		// Load internat struct variables
		row.Init()

		// Run a SELECT query to find out if we need to insert or UPDATE
		selectQuery := fmt.Sprintf(
			"SELECT COUNT(*) FROM %s WHERE %s",
			row.Table,
			row.GetWhere(driver, 0),
		)
		var count int
		err = tx.QueryRow(selectQuery, row.GetPKValues()...).Scan(&count)
		if err != nil {
			tx.Rollback() // rollback the transaction
			return err
		}

		if count == 0 {
			// Primary key not found, let's run an INSERT query
			insertQuery := fmt.Sprintf(
				"INSERT INTO %s(%s) VALUES(%s)",
				row.Table,
				strings.Join(row.GetInsertColumns(), ", "),
				strings.Join(row.GetInsertPlaceholders(driver), ", "),
			)
			_, err := tx.Exec(insertQuery, row.GetInsertValues()...)
			if err != nil {
				tx.Rollback() // rollback the transaction
				return err
			}
		} else {
			// Primary key found, let's run UPDATE query
			updateQuery := fmt.Sprintf(
				"UPDATE %s SET %s WHERE %s",
				row.Table,
				strings.Join(row.GetUpdatePlaceholders(driver), ", "),
				row.GetWhere(driver, row.GetUpdateColumnsLength()),
			)
			values := append(row.GetUpdateValues(), row.GetPKValues()...)
			_, err := tx.Exec(updateQuery, values...)
			if err != nil {
				tx.Rollback() // rollback the transaction
				return err
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // rollback the transaction
		return err
	}

	return nil
}
