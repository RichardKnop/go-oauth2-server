package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetScope() {
	var scope string
	var err error

	// When the requested scope is an empty string,
	// the default scope should be returned
	scope, err = suite.service.getScope("")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar", scope)

	// When the requested scope is valid, it should be returned
	scope, err = suite.service.getScope("foo bar qux")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "foo bar qux", scope)

	// When the requested scope is invalid, an error should be returned
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
