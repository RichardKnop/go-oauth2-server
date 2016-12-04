package oauth_test

import (
	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestFindClientByClientID() {
	var (
		client *models.OauthClient
		err    error
	)

	// When we try to find a client with a bogus client ID
	client, err = suite.service.FindClientByClientID("bogus")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrClientNotFound, err)
	}

	// When we try to find a client with a valid cliend ID
	client, err = suite.service.FindClientByClientID("test_client_1")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client_1", client.Key)
	}
}

func (suite *OauthTestSuite) TestCreateClient() {
	var (
		client *models.OauthClient
		err    error
	)

	// We try to insert a non uniqie client
	client, err = suite.service.CreateClient(
		"test_client_1",           // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrClientIDTaken, err)
	}

	// We try to insert a unique client
	client, err = suite.service.CreateClient(
		"test_client_3",           // client ID
		"test_secret",             // secret
		"https://www.example.com", // redirect URI
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client_3", client.Key)
	}
}

func (suite *OauthTestSuite) TestAuthClient() {
	var (
		client *models.OauthClient
		err    error
	)

	// When we try to authenticate with a bogus client ID
	client, err = suite.service.AuthClient("bogus", "test_secret")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrClientNotFound, err)
	}

	// When we try to authenticate with an invalid secret
	client, err = suite.service.AuthClient("test_client_1", "bogus")

	// Client object should be nil
	assert.Nil(suite.T(), client)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), oauth.ErrInvalidClientSecret, err)
	}

	// When we try to authenticate with valid client ID and secret
	client, err = suite.service.AuthClient("test_client_1", "test_secret")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct client object should be returned
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client_1", client.Key)
	}
}
