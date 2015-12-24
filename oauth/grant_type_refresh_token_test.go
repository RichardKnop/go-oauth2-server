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

func (suite *OauthTestSuite) TestRefreshTokenGrantScopeCannotBeGreater() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "read_write",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read read_write"},
	}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.refreshTokenGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Requested scope cannot be greater\"}",
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrant() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "read_write",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read_write"},
	}

	w := httptest.NewRecorder()
	suite.service.refreshTokenGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.First(accessToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "read_write",
		"refresh_token": "test_token",
	})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
}
