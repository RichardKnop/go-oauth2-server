package oauth_test

import (
	"github.com/RichardKnop/go-oauth2-server/oauth"
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
	_, err = suite.service.GetScope("read_write bogus")
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrInvalidScope, err)
	}
}

func (suite *OauthTestSuite) TestGetDefaultScope() {
	assert.Equal(suite.T(), "read", suite.service.GetDefaultScope())
}

func (suite *OauthTestSuite) TestScopeExists() {
	assert.True(suite.T(), suite.service.ScopeExists("read read_write"))

	assert.False(suite.T(), suite.service.ScopeExists("read_write bogus"))
}
