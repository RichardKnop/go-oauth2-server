package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/response"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthorizationCodeGrantEmptyNotFound() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {""},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	response.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAuthorizationCodeNotFound.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestAuthorizationCodeGrantBogusNotFound() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type": {"authorization_code"},
		"code":       {"bogus"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	response.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAuthorizationCodeNotFound.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestAuthorizationCodeGrantExpired() {
	// Insert a test authorization code
	err := suite.db.Create(&oauth.AuthorizationCode{
		Code:        "test_code",
		ExpiresAt:   time.Now().Add(-10 * time.Second),
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
		"grant_type":   {"authorization_code"},
		"code":         {"test_code"},
		"redirect_uri": {"https://www.example.com"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	response.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAuthorizationCodeExpired.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestAuthorizationCodeGrantInvalidRedirectURI() {
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
		"grant_type":   {"authorization_code"},
		"code":         {"test_code"},
		"redirect_uri": {"https://bogus"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	response.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrInvalidRedirectURI.Error(),
		400,
	)
}

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
		"grant_type":   {"authorization_code"},
		"code":         {"test_code"},
		"redirect_uri": {"https://www.example.com"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Fetch data
	accessToken, refreshToken := new(oauth.AccessToken), new(oauth.RefreshToken)
	assert.False(suite.T(), oauth.AccessTokenPreload(suite.db).
		First(accessToken).RecordNotFound())
	assert.False(suite.T(), oauth.RefreshTokenPreload(suite.db).
		First(refreshToken).RecordNotFound())

	// Check the response
	expected := &oauth.AccessTokenResponse{
		UserID:       accessToken.User.MetaUserID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    oauth.TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	}
	response.TestResponseObject(suite.T(), w, expected, 200)

	// The authorization code should get deleted after use
	assert.True(suite.T(), suite.db.Unscoped().
		First(new(oauth.AuthorizationCode)).RecordNotFound())
}
