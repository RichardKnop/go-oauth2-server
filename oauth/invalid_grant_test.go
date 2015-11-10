package oauth

import (
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *oauthTestSuite) TestInvalidGrantType() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.PostForm = url.Values{"grant_type": {"bogus"}}
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 400, recorded.Recorder.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Invalid grant type\"}",
		recorded.Recorder.Body.String(),
	)
}
