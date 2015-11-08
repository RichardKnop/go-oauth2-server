package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetScope() {
	scope, err := getScope(suite.DB, "")
	assert.Nil(suite.T(), err)
	assert.Equal(
		suite.T(), "foo bar", scope,
		"Should return \"foo bar\"",
	)

	scope, err = getScope(suite.DB, "foo bar qux")
	assert.Nil(suite.T(), err)
	assert.Equal(
		suite.T(), "foo bar qux", scope,
		"Should return \"foo bar qux\"",
	)

	scope, err = getScope(suite.DB, "foo bar bogus")
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid scope", err.Error())
	}
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.Equal(
		suite.T(), true, scopeExists(suite.DB, "foo bar qux"),
		"Should return true",
	)

	assert.Equal(
		suite.T(), false, scopeExists(suite.DB, "foo bar bogus"),
		"Should return false",
	)
}

func TestScopeNotGreater(t *testing.T) {
	assert.Equal(
		t, true, scopeNotGreater("", "foo bar qux"),
		"Should return true",
	)

	assert.Equal(
		t, true, scopeNotGreater("foo", "foo bar qux"),
		"Should return true",
	)

	assert.Equal(
		t, true, scopeNotGreater("foo bar qux", "foo bar qux"),
		"Should return true",
	)

	assert.Equal(
		t, false, scopeNotGreater("foo bar qux bogus", "foo bar qux"),
		"Should return false",
	)
}
