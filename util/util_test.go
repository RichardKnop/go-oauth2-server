package util

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntOrNull(t *testing.T) {
	var nullInt sql.NullInt64
	var value driver.Value
	var err error

	// When the integer is zero
	nullInt = IntOrNull(0)

	// nullInt.Valid should be false
	assert.False(t, nullInt.Valid)

	// nullInt.Value() should return nil
	value, err = nullInt.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the integer is greater than zero
	nullInt = IntOrNull(1)

	// nullInt.Valid should be true
	assert.True(t, nullInt.Valid)

	// nullInt.Value() should return the integer
	value, err = nullInt.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}

func TestStringOrNull(t *testing.T) {
	var nullString sql.NullString
	var value driver.Value
	var err error

	// When the string is empty
	nullString = StringOrNull("")

	// nullString.Valid should be false
	assert.False(t, nullString.Valid)

	// nullString.Value() should return nil
	value, err = nullString.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the string is not empty
	nullString = StringOrNull("foo")

	// nullString.Valid should be true
	assert.True(t, nullString.Valid)

	// nullString.Value() should return the string
	value, err = nullString.Value()
	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("a", []string{"a", "b", "c"}))

	assert.False(t, StringInSlice("d", []string{"a", "b", "c"}))
}

func TestSpaceDelimitedStringNotGreater(t *testing.T) {
	assert.True(t, SpaceDelimitedStringNotGreater("", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("foo", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("bar foo qux", "foo bar qux"))

	assert.False(t, SpaceDelimitedStringNotGreater("foo bar qux bogus", "bar foo qux"))
}
