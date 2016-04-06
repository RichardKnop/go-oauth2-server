package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestRefreshTokenGrantScopeCannotBeGreater() {
	// Insert a test refresh token
	err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Form = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read read_write"},
	}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.refreshTokenGrant(w, r, suite.clients[0])

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", ErrRequestedScopeCannotBeGreater.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantDefaultsToOriginalScope() {
	// Insert a test refresh token
	err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Form = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
	}

	w := httptest.NewRecorder()
	suite.service.refreshTokenGrant(w, r, suite.clients[0])

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&AccessTokenResponse{
		ID:           accessToken.ID,
		UserID:       accessToken.User.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    TokenType,
		Scope:        "read_write",
		RefreshToken: "test_token",
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *OauthTestSuite) TestRefreshTokenGrant() {
	// Insert a test refresh token
	err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Form = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read_write"},
	}

	w := httptest.NewRecorder()
	suite.service.refreshTokenGrant(w, r, suite.clients[0])

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&AccessTokenResponse{
		ID:           accessToken.ID,
		UserID:       accessToken.User.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    TokenType,
		Scope:        "read_write",
		RefreshToken: "test_token",
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
