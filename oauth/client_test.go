package oauth

import (
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestFindClientByClientID() {
	var (
		client *Client
		err    error
	)

	// When we try to find a client with a bogus client ID
	client, err = suite.service.FindClientByClientID("bogus")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errClientNotFound, err)
	}

	// When we try to find a client with a valid cliend ID
	client, err = suite.service.FindClientByClientID("test_client")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.Key)
	}
}

func (suite *OauthTestSuite) TestCreateClient() {
	var (
		client *Client
		err    error
	)

	// We try to insert a non uniqie client
	client, err = suite.service.CreateClient(
		"test_client",             // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "UNIQUE constraint failed: oauth_clients.key", err.Error())
	}

	// We try to insert a unique client
	client, err = suite.service.CreateClient(
		"test_client2",            // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client2", client.Key)
	}
}

func (suite *OauthTestSuite) TestAuthClient() {
	var (
		client *Client
		err    error
	)

	// When we try to authenticate with a bogus client ID
	client, err = suite.service.AuthClient("bogus", "test_secret")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errClientNotFound, err)
	}

	// When we try to authenticate with an invalid secret
	client, err = suite.service.AuthClient("test_client", "bogus")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errInvalidClientSecret, err)
	}

	// When we try to authenticate with valid client ID and secret
	client, err = suite.service.AuthClient("test_client", "test_secret")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.Key)
	}
}
