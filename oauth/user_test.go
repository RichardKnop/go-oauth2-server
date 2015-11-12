package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthUserUsernameNotFound() {
	user, err := suite.service.AuthUser("bogus", "test_password")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUserIncorrectPassword() {
	user, err := suite.service.AuthUser("test@username", "bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid password", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUser() {
	user, err := suite.service.AuthUser("test@username", "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// User should not be nil
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@username", user.Username)
	}
}
