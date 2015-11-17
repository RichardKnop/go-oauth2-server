package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthUser() {
	var user *User
	var err error

	// When we try to authenticate with a bogus username
	user, err = suite.service.AuthUser("bogus", "test_password")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User not found", err.Error())
	}

	// When we try to authenticate with an invalid password
	user, err = suite.service.AuthUser("test@username", "bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid password", err.Error())
	}

	// When we try to authenticate with valid username and password
	user, err = suite.service.AuthUser("test@username", "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@username", user.Username)
	}
}

func (suite *OauthTestSuite) TestFindUserByUsername() {
	var user *User
	var err error

	// When we try to find a user with a bogus username
	user, err = suite.service.FindUserByUsername("bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User not found", err.Error())
	}

	// When we try to find a user with a valid username
	user, err = suite.service.FindUserByUsername("test@username")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@username", user.Username)
	}
}
