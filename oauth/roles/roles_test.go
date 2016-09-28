package roles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGreaterThan(t *testing.T) {
	var (
		is  bool
		err error
	)

	// Superuser is not greater than superuser
	is, err = IsGreaterThan(Superuser, Superuser)
	if assert.Nil(t, err) {
		assert.False(t, is)
	}

	// Superuser is greater than user
	is, err = IsGreaterThan(Superuser, User)
	if assert.Nil(t, err) {
		assert.True(t, is)
	}

	// User is not greater than superuser
	is, err = IsGreaterThan(User, Superuser)
	if assert.Nil(t, err) {
		assert.False(t, is)
	}

	// User is not greater than user
	is, err = IsGreaterThan(User, User)
	if assert.Nil(t, err) {
		assert.False(t, is)
	}
}
