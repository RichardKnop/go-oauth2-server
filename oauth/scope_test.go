package oauth

import (
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

func (suite *OauthTestSuite) TestGetDefaultScope() {
	assert.Equal(suite.T(), "foo bar", s.getDefaultScope())
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.True(suite.T(), suite.service.scopeExists("foo bar qux"))

	assert.False(suite.T(), suite.service.scopeExists("foo bar bogus"))
}
