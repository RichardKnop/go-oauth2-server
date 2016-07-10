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

func (suite *OauthTestSuite) TestPasswordGrant() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"password"},
		"username":   {"test@user"},
		"password":   {"test_password"},
		"scope":      {"read_write"},
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(oauth.AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())
	refreshToken := new(oauth.RefreshToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&oauth.AccessTokenResponse{
		UserID:       accessToken.User.MetaUserID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
