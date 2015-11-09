package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetScope() {
	var scope string
	var err error

	scope, err = getScope(suite.DB, "")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar", scope)

	scope, err = getScope(suite.DB, "foo bar qux")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar qux", scope)

	scope, err = getScope(suite.DB, "foo bar bogus")
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid scope", err.Error())
	}
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.True(suite.T(), scopeExists(suite.DB, "foo bar qux"))
	
	assert.False(suite.T(), scopeExists(suite.DB, "foo bar bogus"))
}

func TestScopeNotGreater(t *testing.T) {
	assert.True(t, scopeNotGreater("", "foo bar qux"))

	assert.True(t, scopeNotGreater("foo", "foo bar qux"))

	assert.True(t, scopeNotGreater("foo bar qux", "foo bar qux"))

	assert.False(t, scopeNotGreater("foo bar qux bogus", "foo bar qux"))
}
