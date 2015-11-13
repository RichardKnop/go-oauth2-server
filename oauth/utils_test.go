package oauth

import (
	"testing"

	"database/sql"
	"database/sql/driver"

	"github.com/stretchr/testify/assert"
)

func TestClientIDOrNull(t *testing.T) {
	var nullInt64 sql.NullInt64
	var value driver.Value
	var err error

	// When client object is nil
	nullInt64 = clientIDOrNull(nil)

	// nullInt64.Valid should be false
	assert.False(t, nullInt64.Valid)

	// nullInt64.Value() should return nil
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When client object is not nil
	nullInt64 = clientIDOrNull(&Client{ID: 1})

	// nullInt64.Valid should be true
	assert.True(t, nullInt64.Valid)

	// nullInt64.Value() should return the object id, in this case int64(1)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}

func TestUserIDOrNull(t *testing.T) {
	var nullInt64 sql.NullInt64
	var value driver.Value
	var err error

	// When user object is nil
	nullInt64 = userIDOrNull(nil)

	// nullInt64.Valid should be false
	assert.False(t, nullInt64.Valid)

	// nullInt64.Value() should return nil
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When user object is not nil
	nullInt64 = userIDOrNull(&User{ID: 1})

	// nullInt64.Valid should be true
	assert.True(t, nullInt64.Valid)

	// nullInt64.Value() should return the object id, in this case int64(1)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}
