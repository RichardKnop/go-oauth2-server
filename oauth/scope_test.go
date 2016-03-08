package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetScope() {
	var (
		scope string
		err   error
	)

	// When the requested scope is an empty string,
	// the default scope should be returned
	scope, err = suite.service.GetScope("")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "read", scope)

	// When the requested scope is valid, it should be returned
	scope, err = suite.service.GetScope("read read_write")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "read read_write", scope)

	// When the requested scope is invalid, an error should be returned
	scope, err = suite.service.GetScope("read_write bogus")
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrInvalidScope, err)
	}
}

func (suite *OauthTestSuite) TestGetDefaultScope() {
	assert.Equal(suite.T(), "read", suite.service.getDefaultScope())
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.True(suite.T(), suite.service.scopeExists("read read_write"))

	assert.False(suite.T(), suite.service.scopeExists("read_write bogus"))
}
