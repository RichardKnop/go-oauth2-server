package oauth2

import (
	"encoding/json"
	"net/url"

	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/stretchr/testify/assert"
)

func (suite *TestSuite) TestClientCredentialsGrant() {
	r := test.MakeSimpleRequest(
		"POST", "http://1.2.3.4/oauth2/api/v1/tokens?grant_type=client_credentials", nil,
	)
	r.SetBasicAuth("test_client_id", "test_client_secret")
	r.PostForm = url.Values{
		"scope": {"bar qux"},
	}
	recorded := test.RunRequest(suite.T(), suite.API.MakeHandler(), r)

	// Status code
	assert.Equal(
		suite.T(),
		200,
		recorded.Recorder.Code, "Status code should be 200",
	)

	// Correct data saved to database
	accessToken := AccessToken{}
	assert.Equal(
		suite.T(),
		false,
		suite.DB.Preload("User").Preload("Client").Preload("RefreshToken").First(&accessToken).RecordNotFound(),
		"Access token should be in the database",
	)
	assert.Equal(
		suite.T(),
		"test_client_id",
		accessToken.Client.ClientID,
		"Access token should belong to test_client_id",
	)
	assert.Equal(
		suite.T(),
		0,
		accessToken.User.ID,
		"Access token should not belong to a user",
	)
	assert.NotEqual(
		suite.T(),
		0,
		accessToken.RefreshToken.ID,
		"Access token should have a refresh token",
	)

	// Response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "bar qux",
		"refresh_token": accessToken.RefreshToken.RefreshToken,
	})
	assert.Equal(
		suite.T(),
		string(expected),
		recorded.Recorder.Body.String(),
		"Response body should be expected access token object",
	)
}
