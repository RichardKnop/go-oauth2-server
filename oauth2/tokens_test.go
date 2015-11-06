package oauth2

import (
	"encoding/json"
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

func (suite *TestSuite) TestTokensPassword() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=password", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
	r.PostForm = url.Values{
		"username": {"test_username"},
		"password": {"test_password"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		200,
		recorded.Recorder.Code, "Status code should be 200",
	)

	accessToken := AccessToken{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.First(&accessToken).RecordNotFound(),
		"Access token should be in the database",
	)

	// RefreshToken record was inserted
	refreshToken := RefreshToken{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.Model(&accessToken).Related(&refreshToken).RecordNotFound(),
		"Refresh token should be in the database",
	)

	// Response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            1,
		"access_token":  accessToken.AccessToken,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "foo bar",
		"refresh_token": refreshToken.RefreshToken,
	})
	assert.Equal(
		suite.T(),
		string(expected),
		recorded.Recorder.Body.String(),
		"Response body should be expected access token object",
	)
}
