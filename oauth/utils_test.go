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

	// Null
	nullInt64 = clientIDOrNull(nil)
	assert.False(t, nullInt64.Valid)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, nil, value)

	// Not null
	nullInt64 = clientIDOrNull(&Client{ID: 1})
	assert.True(t, nullInt64.Valid)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}

func TestUserIDOrNull(t *testing.T) {
	var nullInt64 sql.NullInt64
	var value driver.Value
	var err error

	// Null
	nullInt64 = userIDOrNull(nil)
	assert.False(t, nullInt64.Valid)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, nil, value)

	// Not null
	nullInt64 = userIDOrNull(&User{ID: 1})
	assert.True(t, nullInt64.Valid)
	value, err = nullInt64.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}
