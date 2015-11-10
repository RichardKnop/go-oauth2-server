package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetScope() {
	var scope string
	var err error

	scope, err = suite.service.getScope("")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar", scope)

	scope, err = suite.service.getScope("foo bar qux")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar qux", scope)

	scope, err = suite.service.getScope("foo bar bogus")
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid scope", err.Error())
	}
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.True(suite.T(), suite.service.scopeExists("foo bar qux"))

	assert.False(suite.T(), suite.service.scopeExists("foo bar bogus"))
}

func TestScopeNotGreater(t *testing.T) {
	s := &service{}

	assert.True(t, s.scopeNotGreater("", "foo bar qux"))

	assert.True(t, s.scopeNotGreater("foo", "foo bar qux"))

	assert.True(t, s.scopeNotGreater("foo bar qux", "foo bar qux"))

	assert.False(t, s.scopeNotGreater("foo bar qux bogus", "foo bar qux"))
}
