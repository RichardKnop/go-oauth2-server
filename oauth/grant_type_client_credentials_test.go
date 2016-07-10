package oauth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestClientCredentialsGrant() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"read_write"},
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the access token was inserted
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())

	// Client credentials grant does not produce refresh token
	assert.True(suite.T(), suite.db.Preload("Client").Preload("User").
		First(new(oauth.RefreshToken)).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   3600,
		TokenType:   oauth.TokenType,
		Scope:       "read_write",
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
