package oauth

import (
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *oauthTestSuite) TestAuthClientCredentialsRequired() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client credentials required", err.Error())
	}
}

func (suite *oauthTestSuite) TestAuthClientNotFound() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.SetBasicAuth("bogus", "test_secret")

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *oauthTestSuite) TestAuthClientIncorrectSecret() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.SetBasicAuth("test_client", "bogus")

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *oauthTestSuite) TestAuthClient() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.SetBasicAuth("test_client", "test_secret")

	client, err := suite.service.authClient(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Client should not be nil
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.ClientID)
	}
}

func (suite *oauthTestSuite) TestAuthUserUsernameNotFound() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.PostForm = url.Values{
		"username": {"bogus"},
		"password": {"test_password"},
	}

	user, err := suite.service.authUser(r)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User authentication failed", err.Error())
	}
}

func (suite *oauthTestSuite) TestAuthUserIncorrectPassword() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.PostForm = url.Values{
		"username": {"test_username"},
		"password": {"bogus"},
	}

	user, err := suite.service.authUser(r)

	// User should be nil
	assert.Nil(suite.T(), user)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "User authentication failed", err.Error())
	}
}

func (suite *oauthTestSuite) TestAuthUser() {
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/something", nil)
	r.PostForm = url.Values{
		"username": {"test_username"},
		"password": {"test_password"},
	}

	user, err := suite.service.authUser(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// User should not be nil
	if assert.NotNil(suite.T(), user) {
		assert.Equal(suite.T(), "test_username", user.Username)
	}
}
