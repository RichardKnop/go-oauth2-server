package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthorizationCodeGrantInvalidRedirectURI() {
	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:        "test_code",
		ExpiresAt:   time.Now().Add(+10 * time.Second),
		Client:      suite.client,
		User:        suite.user,
		RedirectURI: util.StringOrNull("https://www.example.com"),
		Scope:       "foo",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Prepare a request object
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_code"},
	}

	w := httptest.NewRecorder()
	suite.service.authorizationCodeGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Invalid redirect URI\"}",
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestAuthorizationCodeGrant() {
	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:        "test_code",
		ExpiresAt:   time.Now().Add(+10 * time.Second),
		Client:      suite.client,
		User:        suite.user,
		RedirectURI: util.StringOrNull("https://www.example.com"),
		Scope:       "foo",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Prepare a request object
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {"test_code"},
		"redirect_uri": {"https://www.example.com"},
	}

	w := httptest.NewRecorder()
	suite.service.authorizationCodeGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.First(accessToken).RecordNotFound())
	refreshToken := new(RefreshToken)
	assert.False(suite.T(), suite.db.First(refreshToken).RecordNotFound())

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "foo",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))

	// Check the authorization code was deleted
	assert.True(suite.T(), suite.db.First(new(AuthorizationCode)).RecordNotFound())
}
