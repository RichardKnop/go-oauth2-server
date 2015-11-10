package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthorizationCodeGrantNotFound() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_auth_code"},
	}

	w := httptest.NewRecorder()
	suite.service.authorizationCodeGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	expected := "{\"error\":\"Authorization code not found\"}"
	assert.Equal(suite.T(), expected, strings.TrimSpace(w.Body.String()))
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
	r, err := http.NewRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_auth_code"},
	}

	w := httptest.NewRecorder()
	suite.service.authorizationCodeGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

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
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))

	// Check the authorization code was deleted
	assert.True(suite.T(), suite.db.First(&AuthorizationCode{}).RecordNotFound())
}
