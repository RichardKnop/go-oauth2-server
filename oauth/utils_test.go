package oauth

import (
	"testing"

	"database/sql"

	"github.com/stretchr/testify/assert"
)

func TestClientIDOrNull(t *testing.T) {
	var nullInt64 sql.NullInt64

	nullInt64 = clientIDOrNull(nil)
	assert.False(t, nullInt64.Valid)

	nullInt64 = clientIDOrNull(&Client{ID: 1})
	assert.True(t, nullInt64.Valid)
}

func TestUserIDOrNull(t *testing.T) {
	var nullInt64 sql.NullInt64

	nullInt64 = userIDOrNull(nil)
	assert.False(t, nullInt64.Valid)

	nullInt64 = userIDOrNull(&User{ID: 1})
	assert.True(t, nullInt64.Valid)
}
