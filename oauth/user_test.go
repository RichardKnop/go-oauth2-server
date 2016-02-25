package oauth

import (
	"github.com/stretchr/testify/assert"

	pass "github.com/RichardKnop/go-oauth2-server/password"
	"github.com/RichardKnop/go-oauth2-server/util"
)

func (suite *OauthTestSuite) TestFindUserByUsername() {
	var (
		user *User
		err  error
	)

	// When we try to find a user with a bogus username
	user, err = suite.service.FindUserByUsername("bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserNotFound, err)
	}

	// When we try to find a user with a valid username
	user, err = suite.service.FindUserByUsername("test@user")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.Username)
	}
}

func (suite *OauthTestSuite) TestCreateUser() {
	var (
		user *User
		err  error
	)

	// We try to insert a non unique user
	user, err = suite.service.CreateUser(
		"test@user",     // username
		"test_password", // password
	)

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_users.username", err.Error())
	}

	// We try to insert a unique user
	user, err = suite.service.CreateUser(
		"test@user2",    // username
		"test_password", // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user2", user.Username)
	}
}

func (suite *OauthTestSuite) TestSetPassword() {
	var (
		user *User
		err  error
	)

	// Insert a test user without a password
	user = &User{
		Username: "test@user_nopass",
		Password: util.StringOrNull(""),
	}
	err = suite.db.Create(user).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Try to set an empty password
	err = suite.service.SetPassword(user, "")

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errCannotSetEmptyUserPassword, err)
	}

	// Try changing the password
	err = suite.service.SetPassword(user, "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// User object should have been updated
	assert.Equal(suite.T(), "test@user_nopass", user.Username)
	assert.Nil(suite.T(), pass.VerifyPassword(user.Password.String, "test_password"))
}

func (suite *OauthTestSuite) TestAuthUser() {
	var (
		user *User
		err  error
	)

	// Insert a test user without a password
	err = suite.db.Create(&User{
		Username: "test@user_nopass",
		Password: util.StringOrNull(""),
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// When we try to authenticate a user without a password
	user, err = suite.service.AuthUser("test@user_nopass", "bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserPasswordNotSet, err)
	}

	// When we try to authenticate with a bogus username
	user, err = suite.service.AuthUser("bogus", "test_password")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserNotFound, err)
	}

	// When we try to authenticate with an invalid password
	user, err = suite.service.AuthUser("test@user", "bogus")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errInvalidUserPassword, err)
	}

	// When we try to authenticate with valid username and password
	user, err = suite.service.AuthUser("test@user", "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user", user.Username)
	}
}

func (suite *OauthTestSuite) TestBlankPassword() {
	var (
		user *User
		err  error
	)

	user, err = suite.service.CreateUser(
		"test@user_nopass", // username
		"",                 // password
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct user object should be returned
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test@user_nopass", user.Username)
	}

	// When we try to authenticate
	user, err = suite.service.AuthUser("test@user_nopass", "")

	// User object should be nil
	assert.Nil(suite.T(), user)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errUserPasswordNotSet, err)
	}
}
