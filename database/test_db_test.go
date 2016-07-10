package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDBName = "area_database_test"
	testDBUser = "area"
)

func TestCreateTestDatabaseFailsWithBadValues(t *testing.T) {
	db, err := CreateTestDatabase("!_@£@$@!±/\\", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateTestDatabaseWorksWithValidEntry(t *testing.T) {
	db, err := CreateTestDatabase("", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Close()
	assert.Nil(t, err)
}

func TestCreateTestDatabaseFailsWithMissingFixtureFile(t *testing.T) {
	badFixtures := []string{"/badfilename"}
	db, err := CreateTestDatabase("", nil, badFixtures)
	assert.EqualError(t, err, "Error loading file /badfilename: open /badfilename: no such file or directory")
	assert.Nil(t, db)
}

func TestCreateTestDatabasePostgresFailsWithBadValues(t *testing.T) {
	db, err := CreateTestDatabasePostgres("", "", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateTestDatabasePostgresWorksWithValidEntry(t *testing.T) {
	db, err := CreateTestDatabasePostgres(testDBUser, testDBName, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = db.Close()
	assert.Nil(t, err)
}

func TestCreateTestDatabasePostgresFailsWithMissingFixtureFile(t *testing.T) {
	badFixtures := []string{"/badfilename"}
	db, err := CreateTestDatabasePostgres(testDBUser, testDBName, nil, badFixtures)
	assert.EqualError(t, err, "Error loading file /badfilename: open /badfilename: no such file or directory")
	assert.Nil(t, db)
}

func TestOpenPostgresDBFailsWhenDBNotFound(t *testing.T) {
	db, err := openPostgresDB(testDBUser, "missingDB")
	assert.EqualError(t, err, "pq: database \"missingDB\" does not exist")
	assert.Nil(t, db)
}

func TestOpenPostgresDBFailsWhenUserNotFound(t *testing.T) {
	db, err := openPostgresDB("bogus_user", testDBName)
	assert.EqualError(t, err, "pq: role \"bogus_user\" does not exist")
	assert.Nil(t, db)
}
