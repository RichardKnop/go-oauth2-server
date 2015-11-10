package oauth

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthorizationCodeGrantNotFound() {
	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_auth_code"},
	}
	recorded := test.RunRequest(suite.T(), suite.api.MakeHandler(), r)

	// Check the status code
	assert.Equal(suite.T(), 400, recorded.Recorder.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Authorization code not found\"}",
		recorded.Recorder.Body.String(),
	)
}

func (suite *OauthTestSuite) TestAuthorizationCodeGrant() {
	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:      "test_auth_code",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.client,
		Scope:     "foo bar",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r := test.MakeSimpleRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_auth_code"},
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
		"scope":         "foo bar",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(suite.T(), string(expected), recorded.Recorder.Body.String())

	// Check the authorization code was deleted
	assert.True(suite.T(), suite.db.First(&AuthorizationCode{}).RecordNotFound())
}
