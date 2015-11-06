package oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestScopeExists() {
	assert.Equal(
		suite.T(),
		true,
		scopeExists(suite.DB, "foo bar qux"), "Should return true",
	)

	assert.Equal(
		suite.T(),
		false,
		scopeExists(suite.DB, "foo bar bogus"), "Should return false",
	)
}

func TestScopeNotGreater(t *testing.T) {
	assert.Equal(
		t,
		true,
		scopeNotGreater("foo", "foo bar qux"), "Should return true",
	)

	assert.Equal(
		t,
		true,
		scopeNotGreater("foo bar qux", "foo bar qux"), "Should return true",
	)

	assert.Equal(
		t,
		false,
		scopeNotGreater("foo bar qux bogus", "foo bar qux"), "Should return false",
	)
}
