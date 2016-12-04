package password_test

import (
	"testing"

	"github.com/RichardKnop/go-oauth2-server/util/password"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	// Test valid passwords
	assert.Nil(t, password.VerifyPassword(
		"$2a$10$CUoGytf1pR7CC6Y043gt/.vFJUV4IRqvH5R6F0VfITP8s2TqrQ.4e",
		"test_secret",
	))

	assert.Nil(t, password.VerifyPassword(
		"$2a$10$4J4t9xuWhOKhfjN0bOKNReS9sL3BVSN9zxIr2.VaWWQfRBWh1dQIS",
		"test_password",
	))

	// Test invalid password
	assert.NotNil(t, password.VerifyPassword("bogus", "password"))
}
