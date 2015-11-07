package oauth

import (
	"log"
	"net/url"
	"time"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestRefreshTokenGrantNotFound() {
	// Make a request
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
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
		RefreshToken: "test_refresh_token",
		ExpiresAt:    time.Now().Add(-1 * time.Second),
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
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

func (suite *OauthTestSuite) TestRefreshTokenGrantAccessTokenNotFound() {
	// Insert a test refresh token
	if err := suite.DB.Create(&RefreshToken{
		RefreshToken: "test_refresh_token",
		ExpiresAt:    time.Now().Add(+1 * time.Second),
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
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
		"{\"error\":\"Access token not found\"}",
		recorded.Recorder.Body.String(),
		"Body should be expected JSON error",
	)
}
