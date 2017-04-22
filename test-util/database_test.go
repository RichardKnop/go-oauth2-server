package testutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/RichardKnop/go-oauth2-server/test-util"
)

var (
	testDBName = "go_oauth2_server_database_test"
	testDBUser = "go_oauth2_server"
)

func TestCreateTestDatabaseFailsWithBadValues(t *testing.T) {
	db, err := testutil.CreateTestDatabase("!_@£@$@!±/\\", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateTestDatabaseWorksWithValidEntry(t *testing.T) {
	db, err := testutil.CreateTestDatabase("", nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	err = db.Close()
	assert.Nil(t, err)
}

func TestCreateTestDatabaseFailsWithMissingFixtureFile(t *testing.T) {
	badFixtures := []string{"/badfilename"}
	db, err := testutil.CreateTestDatabase("", nil, badFixtures)
	assert.EqualError(t, err, "Error loading file /badfilename: open /badfilename: no such file or directory")
	assert.Nil(t, db)
}

func TestCreateTestDatabasePostgresFailsWithBadValues(t *testing.T) {
	db, err := testutil.CreateTestDatabasePostgres("", "", nil, nil)
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCreateTestDatabasePostgresWorksWithValidEntry(t *testing.T) {
	db, err := testutil.CreateTestDatabasePostgres(testDBUser, testDBName, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	err = db.Close()
	assert.Nil(t, err)
}

func TestCreateTestDatabasePostgresFailsWithMissingFixtureFile(t *testing.T) {
	badFixtures := []string{"/badfilename"}
	db, err := testutil.CreateTestDatabasePostgres(testDBUser, testDBName, nil, badFixtures)
	assert.EqualError(t, err, "Error loading file /badfilename: open /badfilename: no such file or directory")
	assert.Nil(t, db)
}
