package oauth

import (
	"log"
	"net/http"
	"net/url"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthClientCredentialsRequired() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client credentials required", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClientNotFound() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth("bogus", "test_secret")

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClientIncorrectSecret() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth("test_client", "bogus")

	client, err := suite.service.authClient(r)

	// Client should be nil
	assert.Nil(suite.T(), client)

	// Error should not be nil
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthClient() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth("test_client", "test_secret")

	client, err := suite.service.authClient(r)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Client should not be nil
	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client", client.ClientID)
	}
}

func (suite *OauthTestSuite) TestAuthUserUsernameNotFound() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
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

func (suite *OauthTestSuite) TestAuthUserIncorrectPassword() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
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

func (suite *OauthTestSuite) TestAuthUser() {
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
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
