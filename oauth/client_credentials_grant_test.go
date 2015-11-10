package oauth

import (
	"encoding/json"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestClientCredentialsGrant() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"bar qux"},
	}
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 200, recorded.Recorder.Code)

	// Check the correct data was inserted
	accessToken := AccessToken{}
	assert.False(suite.T(), suite.db.First(&accessToken).RecordNotFound())
	refreshToken := RefreshToken{}
	assert.False(suite.T(), suite.db.First(&refreshToken).RecordNotFound())

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "bar qux",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(suite.T(), string(expected), recorded.Recorder.Body.String())
}
