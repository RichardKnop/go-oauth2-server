package fixtures_test

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/RichardKnop/go-fixtures"
	"github.com/stretchr/testify/assert"
	// Driver
	_ "github.com/lib/pq"
)

const testPostgresDbUser = "go_fixtures"

func TestLoadWorksWithValidDataPostgres(t *testing.T) {
	t.Parallel()

	var (
		db  *sql.DB
		err error
	)

	// Connect to a test Postgres db
	db, err = rebuildDatabasePostgres(testPostgresDbUser, "go_fixtures_test_load")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a test schema
	_, err = db.Exec(testSchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	// Let's load the fixture, since the database is empty, this should run inserts
	err = fixtures.Load([]byte(testData), db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	var (
		count        int
		rows         *sql.Rows
		id           int
		stringField  string
		booleanField bool
		intField     int
		createdAt    *time.Time
		updatedAt    *time.Time
		someID       int
		otherID      int
	)

	// Check row counts
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 1, count)

	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 1, count)

	// Check correct data has been loaded into some_table
	rows, err = db.Query("SELECT id, string_field, boolean_field, " +
		"created_at, updated_at FROM some_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&stringField,
			&booleanField,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 1, id)
		assert.Equal(t, "foobar", stringField)
		assert.Equal(t, true, booleanField)
		assert.NotNil(t, createdAt)
		assert.Nil(t, updatedAt)
	}

	// Check correct data has been loaded into other_table
	rows, err = db.Query("SELECT id, int_field, boolean_field, " +
		"created_at, updated_at FROM other_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&intField,
			&booleanField,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 2, id)
		assert.Equal(t, 123, intField)
		assert.Equal(t, false, booleanField)
		assert.NotNil(t, createdAt)
		assert.Nil(t, updatedAt)
	}

	// Check correct data has been loaded into join_table
	rows, err = db.Query("SELECT some_id, other_id FROM join_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&someID,
			&otherID,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 1, someID)
		assert.Equal(t, 2, otherID)
	}

	// Let's reload the fixture, this should run updates
	err = fixtures.Load([]byte(testData), db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	// Check row counts, should be unchanged
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 1, count)

	// Check correct data has been loaded into some_table
	rows, err = db.Query("SELECT id, string_field, boolean_field, " +
		"created_at, updated_at FROM some_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&stringField,
			&booleanField,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 1, id)
		assert.Equal(t, "foobar", stringField)
		assert.Equal(t, true, booleanField)
		assert.NotNil(t, createdAt)
		assert.NotNil(t, updatedAt)
	}

	// Check correct data has been loaded into other_table
	rows, err = db.Query("SELECT id, int_field, boolean_field, " +
		"created_at, updated_at FROM other_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&intField,
			&booleanField,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 2, id)
		assert.Equal(t, 123, intField)
		assert.Equal(t, false, booleanField)
		assert.NotNil(t, createdAt)
		assert.NotNil(t, updatedAt)
	}

	// Check correct data has been loaded into join_table
	rows, err = db.Query("SELECT some_id, other_id FROM join_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&someID,
			&otherID,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 1, someID)
		assert.Equal(t, 2, otherID)
	}
}

func TestLoadFileWorksWithValidFilePostgres(t *testing.T) {
	t.Parallel()

	var (
		db  *sql.DB
		err error
	)

	// Connect to a test Postgres db
	db, err = rebuildDatabasePostgres(testPostgresDbUser, "go_fixtures_test_load_file")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a test schema
	_, err = db.Exec(testSchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	// Check row counts to show no data
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 0, count)

	// Let's load the fixture, since the database is empty, this should run inserts
	err = fixtures.LoadFile(fixtureFile, db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	var (
		rows         *sql.Rows
		id           int
		stringField  string
		booleanField bool
		createdAt    *time.Time
		updatedAt    *time.Time
	)

	// Check row counts
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 0, count)

	// Check correct data has been loaded into some_table
	rows, err = db.Query("SELECT id, string_field, boolean_field, " +
		"created_at, updated_at FROM some_table")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&id,
			&stringField,
			&booleanField,
			&createdAt,
			&updatedAt,
		); err != nil {
			log.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 1, id)
		assert.Equal(t, "foobar", stringField)
		assert.Equal(t, true, booleanField)
		assert.NotNil(t, createdAt)
		assert.Nil(t, updatedAt)
	}

	// Let's reload the fixture, this should run updates
	err = fixtures.LoadFile(fixtureFile, db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	// Check row counts, should be unchanged
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 0, count)
}

func TestLoadFileFailsWithMissingFilePostgres(t *testing.T) {
	t.Parallel()

	var (
		db  *sql.DB
		err error
	)

	// Connect to a test Postgres db
	db, err = rebuildDatabasePostgres(testPostgresDbUser, "go_fixtures_test_load_file_missing_file")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a test schema
	_, err = db.Exec(testSchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	// Let's load the fixture, since the database is empty, this should run inserts
	err = fixtures.LoadFile("bad_filename.yml", db, "postgres")

	// Error should be nil
	assert.EqualError(t, err, "Error loading file bad_filename.yml: open bad_filename.yml: no such file or directory")
}

func TestLoadFilesWorksWithValidFilesPostgres(t *testing.T) {
	t.Parallel()

	var (
		db  *sql.DB
		err error
	)

	// Connect to a test Postgres db
	db, err = rebuildDatabasePostgres(testPostgresDbUser, "go_fixtures_test_load_files")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a test schema
	_, err = db.Exec(testSchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	var count int

	// Check rows are empty first
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 0, count)

	// Let's load the fixture, since the database is empty, this should run inserts
	err = fixtures.LoadFiles(fixtureFiles, db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	// Check row counts
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 1, count)

	// Let's reload the fixtures, this should run updates
	err = fixtures.LoadFiles(fixtureFiles, db, "postgres")

	// Error should be nil
	assert.Nil(t, err)

	// Check row counts, should be unchanged
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 1, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 1, count)
}

func TestLoadFilesFailsWithABadFilePostgres(t *testing.T) {
	t.Parallel()

	var (
		db  *sql.DB
		err error
	)

	// Connect to a test Postgres db
	db, err = rebuildDatabasePostgres(testPostgresDbUser, "go_fixtures_test_load_files_bad_file")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a test schema
	_, err = db.Exec(testSchemaPostgres)
	if err != nil {
		log.Fatal(err)
	}

	var count int

	// Check rows are empty first
	db.QueryRow("SELECT COUNT(*) FROM some_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM other_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM join_table").Scan(&count)
	assert.Equal(t, 0, count)
	db.QueryRow("SELECT COUNT(*) FROM string_key_table").Scan(&count)
	assert.Equal(t, 0, count)

	var badList = []string{
		fixtureFile,
		"bad_file",
	}

	// Let's load the fixture, since the database is empty, this should run inserts
	err = fixtures.LoadFiles(badList, db, "postgres")

	// Error should be nil
	assert.EqualError(t, err, "Error loading file bad_file: open bad_file: no such file or directory")
}

// rebuildDatabase attempts to delete an existing Postgres
// database and rebuild it, returning a pointer to it
func rebuildDatabasePostgres(dbUser, dbName string) (*sql.DB, error) {

	dropPostgresDB(dbUser, dbName)

	if err := createPostgresDB(dbUser, dbName); err != nil {
		return nil, err
	}

	return openPostgresDB(dbUser, dbName)
}

func openPostgresDB(dbUser, dbName string) (*sql.DB, error) {
	// Init a new postgres test database connection
	return sql.Open("postgres",
		fmt.Sprintf(
			"sslmode=disable host=localhost port=5432 user=%s password='' dbname=%s",
			dbUser,
			dbName,
		),
	)
}

func createPostgresDB(dbUser, dbName string) error {
	// Create a new test database
	createDbCmd := fmt.Sprintf("createdb -U %s %s", dbUser, dbName)
	log.Println(createDbCmd)
	out, err := exec.Command("sh", "-c", createDbCmd).Output()
	if err != nil {
		log.Printf("%v", string(out))
		return err
	}
	return nil
}

func dropPostgresDB(dbUser, dbName string) {
	// Delete the current database if it exists
	dropDbCmd := fmt.Sprintf("dropdb --if-exists -U %s %s", dbUser, dbName)
	fmt.Println(dropDbCmd)
	exec.Command("sh", "-c", dropDbCmd).Output()
}
