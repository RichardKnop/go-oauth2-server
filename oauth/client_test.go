package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthClientNotFound() {
	// When we try to authenticate with bogus client id
	client, err := suite.service.AuthClient("bogus", "test_secret")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClientIncorrectSecret() {
	// When we try to authenticate with invalid secret
	client, err := suite.service.AuthClient("test_client", "bogus")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Invalid secret", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClient() {
	// When we try to authenticate with valid client id and secret
	client, err := suite.service.AuthClient("test_client", "test_secret")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.ClientID)
	}
}
