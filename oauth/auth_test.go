package oauth

import "github.com/stretchr/testify/assert"

func (suite *OauthTestSuite) TestAuthClientNotFound() {
	client, err := suite.service.AuthClient("bogus", "test_secret")

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClientIncorrectSecret() {
	client, err := suite.service.AuthClient("test_client", "bogus")

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClient() {
	client, err := suite.service.AuthClient("test_client", "test_secret")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Client should not be nil
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.ClientID)
	}
}

func (suite *OauthTestSuite) TestAuthUserUsernameNotFound() {
	user, err := suite.service.AuthUser("bogus", "test_password")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUserIncorrectPassword() {
	user, err := suite.service.AuthUser("test_username", "bogus")

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthUser() {
	user, err := suite.service.AuthUser("test_username", "test_password")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// User should not be nil
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_username", user.Username)
	}
}
