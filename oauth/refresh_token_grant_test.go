package oauth

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestRefreshTokenGrantNotFound() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"bogus"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Check the status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Check the response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Refresh token not found\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantExpired() {
	// Insert a test refresh token
	if err := suite.DB.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_refresh_token"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Check the status code
	assert.Equal(
		suite.T(),
		400,
		recorded.Recorder.Code, "Status code should be 400",
	)

	// Check the response body
	assert.Equal(
		suite.T(),
		"{\"error\":\"Refresh token expired\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrant() {
	// Insert a test refresh token
	if err := suite.DB.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "foo bar",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_refresh_token"},
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

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "foo bar",
		"refresh_token": "test_refresh_token",
	})
	assert.Equal(
		suite.T(),
		string(expected),
		recorded.Recorder.Body.String(),
		"Response body should be expected access token object",
	)
}
