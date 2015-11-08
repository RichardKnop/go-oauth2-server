package oauth

import (
	"encoding/json"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestPasswordGrant() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"password"},
		"username":   {"test_username"},
		"password":   {"test_password"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Check the status code
	assert.Equal(
		suite.T(),
		200,
		recorded.Recorder.Code, "Status code should be 200",
	)

	// Check the correct data was inserted
	accessToken := AccessToken{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.First(&accessToken).RecordNotFound(),
		"Access token should be in the database",
	)
	refreshToken := RefreshToken{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.First(&refreshToken).RecordNotFound(),
		"Refresh token should be in the database",
	)

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "foo bar",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(
		suite.T(),
		string(expected),
		recorded.Recorder.Body.String(),
		"Response body should be expected access token object",
	)
}
