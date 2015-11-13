package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthUserUsernameNotFound() {
	// When we try to authenticate with bogus username
	user, err := suite.service.AuthUser("bogus", "test_password")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUserIncorrectPassword() {
	// When we try to authenticate with invalid password
	user, err := suite.service.AuthUser("test@username", "bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid password", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUser() {
	// When we try to authenticate with valid username and password
	user, err := suite.service.AuthUser("test@username", "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@username", user.Username)
	}
}
