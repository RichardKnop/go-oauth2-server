package oauth2

import (
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
