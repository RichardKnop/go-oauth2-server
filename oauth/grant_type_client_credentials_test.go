package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/response"
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

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Fetch data
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), oauth.AccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())

	// Check the response
	expected := &oauth.AccessTokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   3600,
		TokenType:   oauth.TokenType,
		Scope:       "read_write",
	}
	response.TestResponseObject(suite.T(), w, expected, 200)

	// Client credentials grant does not produce refresh token
	assert.True(suite.T(), oauth.RefreshTokenPreload(suite.db).
		First(new(oauth.RefreshToken)).RecordNotFound())
}
