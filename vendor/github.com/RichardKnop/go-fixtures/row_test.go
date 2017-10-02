package fixtures_test

import (
	"testing"

	"github.com/RichardKnop/go-fixtures"
	"github.com/stretchr/testify/assert"
)

func TestRow(t *testing.T) {
	t.Parallel()

	// Create a test Row instance
	row := &fixtures.Row{
		Table: "some_table",
		PK: map[string]interface{}{
			"some_id":  interface{}(1),
			"other_id": interface{}(2),
		},
		Fields: map[string]interface{}{
			"string_field":  interface{}("foobar"),
			"boolean_field": interface{}(true),
			"created_at":    interface{}("ON_INSERT_NOW()"),
			"updated_at":    interface{}("ON_UPDATE_NOW()"),
		},
	}

	// Run the init method to populate internal variables
	row.Init()

	var (
		expectedString     string
		expectedStrings    []string
		expectedInterfaces []interface{}
	)

	// Test insert and update column lengths
	assert.Equal(t, 5, row.GetInsertColumnsLength())
	assert.Equal(t, 5, row.GetUpdateColumnsLength())

	// Test insert and update columns
	expectedStrings = []string{"\"other_id\"", "\"some_id\"",
		"\"boolean_field\"", "\"created_at\"", "\"string_field\""}
	assert.Equal(t, expectedStrings, row.GetInsertColumns())
	expectedStrings = []string{"\"other_id\"", "\"some_id\"",
		"\"boolean_field\"", "\"string_field\"", "\"updated_at\""}
	assert.Equal(t, expectedStrings, row.GetUpdateColumns())

	// Test postgres placeholders ($1, $2 and so on)
	expectedStrings = []string{"$1", "$2", "$3", "$4", "$5"}
	assert.Equal(t, expectedStrings, row.GetInsertPlaceholders("postgres"))
	expectedStrings = []string{"\"other_id\" = $1", "\"some_id\" = $2",
		"\"boolean_field\" = $3", "\"string_field\" = $4", "\"updated_at\" = $5"}
	assert.Equal(t, expectedStrings, row.GetUpdatePlaceholders("postgres"))

	// Test non postgres placeholders (?)
	expectedStrings = []string{"?", "?", "?", "?", "?"}
	assert.Equal(t, expectedStrings, row.GetInsertPlaceholders("sqlite"))
	expectedStrings = []string{"\"other_id\" = ?", "\"some_id\" = ?",
		"\"boolean_field\" = ?", "\"string_field\" = ?", "\"updated_at\" = ?"}
	assert.Equal(t, expectedStrings, row.GetUpdatePlaceholders("sqlite"))

	// Test where clause
	expectedString = "other_id = $3 AND some_id = $4"
	assert.Equal(t, expectedString, row.GetWhere("postgres", 2))

	// Test primary key values
	expectedInterfaces = []interface{}{interface{}(2), interface{}(1)}
	assert.Equal(t, expectedInterfaces, row.GetPKValues())
}
