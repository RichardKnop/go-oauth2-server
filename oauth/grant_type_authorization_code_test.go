package oauth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthorizationCodeGrant() {
	// Insert a test authorization code
	err := suite.db.Create(&oauth.AuthorizationCode{
		Code:        "test_code",
		ExpiresAt:   time.Now().Add(+10 * time.Second),
		Client:      suite.clients[0],
		User:        suite.users[0],
		RedirectURI: util.StringOrNull("https://www.example.com"),
		Scope:       "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"test_code"},
	}

	// First we will test an invalid redirect URI error
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrInvalidRedirectURI.Error()),
		strings.TrimSpace(w.Body.String()),
	)

	// Now add the redirect URI parameter
	r.Form.Set("redirect_uri", "https://www.example.com")

	// And test a successful case
	w = httptest.NewRecorder()
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

	// Check the authorization code was deleted
	notFound := suite.db.Unscoped().First(new(oauth.AuthorizationCode)).RecordNotFound()
	assert.True(suite.T(), notFound)
}
