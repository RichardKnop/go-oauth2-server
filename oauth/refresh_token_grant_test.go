package oauth

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *oauthTestSuite) TestRefreshTokenGrantNotFound() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"bogus"},
	}
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 400, recorded.Recorder.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Refresh token not found\"}",
		recorded.Recorder.Body.String(),
	)
}

func (suite *oauthTestSuite) TestRefreshTokenGrantExpired() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.client,
		User:      *suite.user,
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
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 400, recorded.Recorder.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Refresh token expired\"}",
		recorded.Recorder.Body.String(),
	)
}

func (suite *oauthTestSuite) TestRefreshTokenGrant() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.client,
		User:      *suite.user,
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
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 200, recorded.Recorder.Code)

	// Check the correct data was inserted
	accessToken := AccessToken{}
	assert.False(suite.T(), suite.db.First(&accessToken).RecordNotFound())

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "foo bar",
		"refresh_token": "test_refresh_token",
	})
	assert.Equal(suite.T(), string(expected), recorded.Recorder.Body.String())
}
