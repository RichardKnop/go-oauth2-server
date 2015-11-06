package oauth2

import (
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestAuthClientAuthenticationRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)

	client, err := authClient(r, suite.DB)

	assert.Nil(suite.T(), client)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication required", err.Error())
	}
}

func (suite *TestSuite) TestAuthClientNotFound() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("bogus", "test_client_secret")

	client, err := authClient(r, suite.DB)

	assert.Nil(suite.T(), client)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *TestSuite) TestAuthClientIncorrectSecret() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("test_client_id", "bogus")

	client, err := authClient(r, suite.DB)

	assert.Nil(suite.T(), client)

	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Client authentication failed", err.Error())
	}
}

func (suite *TestSuite) TestAuthClient() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")

	client, err := authClient(r, suite.DB)

	assert.Nil(suite.T(), err)

	if assert.NotNil(suite.T(), client) {
		assert.Equal(suite.T(), "test_client_id", client.ClientID)
	}
}
