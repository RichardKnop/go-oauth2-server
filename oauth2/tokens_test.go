package oauth2

import (
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestTokensInvalidGrantType() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=bogus", nil,
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Invalid grant type\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *TestSuite) TestTokensPasswordUserClientAuthenticationRequired() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		401,
		recorded.Recorder.Code, "Status code should be 401",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Client authentication required\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)

	// TODO - WWW-Authenticate header
}

func (suite *TestSuite) TestTokensPasswordClientNotFound() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("bogus", "test_client_secret")
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		401,
		recorded.Recorder.Code, "Status code should be 401",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Client authentication failed\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)

	// TODO - WWW-Authenticate header
}

func (suite *TestSuite) TestTokensPasswordIncorrectClientSecret() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("test_client_id", "bogus")
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		401,
		recorded.Recorder.Code, "Status code should be 401",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Client authentication failed\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)

	// TODO - WWW-Authenticate header
}

func (suite *TestSuite) TestTokensPasswordUserNotFound() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
	r.PostForm = url.Values{
		"username": {"bogus"},
		"password": {"test_password"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		401,
		recorded.Recorder.Code, "Status code should be 401",
	)

	// Response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"User authentication failed\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)

	// TODO - WWW-Authenticate header
}
