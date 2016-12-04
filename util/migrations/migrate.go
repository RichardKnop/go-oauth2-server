package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// MigrationStage ...
type MigrationStage struct {
	Name     string
	Function func(db *gorm.DB, name string) error
}

// Migrate ...
func Migrate(db *gorm.DB, migrations []MigrationStage) error {
	for _, m := range migrations {
		if MigrationExists(db, m.Name) {
			continue
		}

		if err := m.Function(db, m.Name); err != nil {
			return err
		}

		if err := SaveMigration(db, m.Name); err != nil {
			return err
		}
	}

	return nil
}

// MigrateAll runs bootstrap, then all migration functions listed against
// the specified database and logs any errors
func MigrateAll(db *gorm.DB, migrationFunctions []func(*gorm.DB) error) {
	if err := Bootstrap(db); err != nil {
		logger.Error(err)
	}

	for _, m := range migrationFunctions {
		if err := m(db); err != nil {
			logger.Error(err)
		}
	}
}

// MigrationExists checks if the migration called migrationName has been run already
func MigrationExists(db *gorm.DB, migrationName string) bool {
	migration := new(Migration)
	found := !db.Where("name = ?", migrationName).First(migration).RecordNotFound()

	if found {
		logger.Infof("Skipping %s migration", migrationName)
	} else {
		logger.Infof("Running %s migration", migrationName)
	}

	return found
}

// SaveMigration saves a migration to the migration table
func SaveMigration(db *gorm.DB, migrationName string) error {
	migration := new(Migration)
	migration.Name = migrationName

	if err := db.Create(migration).Error; err != nil {
		logger.Errorf("Error saving record to migrations table: %s", err)
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
