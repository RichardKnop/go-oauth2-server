package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/oauth/tokentypes"
	"github.com/RichardKnop/go-oauth2-server/test-util"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestRefreshTokenGrantEmptyNotFound() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {""},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRefreshTokenNotFound.Error(),
		404,
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantBogusNotFound() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"bogus_token"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRefreshTokenNotFound.Error(),
		404,
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantExipired() {
	// Insert a test refresh token
	err := suite.db.Create(&models.OauthRefreshToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token",
		ExpiresAt: time.Now().UTC().Add(-10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read read_write"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRefreshTokenExpired.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantScopeCannotBeGreater() {
	// Insert a test refresh token
	err := suite.db.Create(&models.OauthRefreshToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token",
		ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read read_write"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRequestedScopeCannotBeGreater.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestRefreshTokenGrantDefaultsToOriginalScope() {
	// Insert a test refresh token
	err := suite.db.Create(&models.OauthRefreshToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token",
		ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Fetch data
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		Last(accessToken).RecordNotFound())

	// Check the response body
	expected := &oauth.AccessTokenResponse{
		UserID:       accessToken.UserID.String,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: "test_token",
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}

func (suite *OauthTestSuite) TestRefreshTokenGrant() {
	// Insert a test refresh token
	err := suite.db.Create(&models.OauthRefreshToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token",
		ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {"test_token"},
		"scope":         {"read_write"},
	}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Fetch data
	accessToken := new(models.OauthAccessToken)
	assert.False(suite.T(), models.OauthAccessTokenPreload(suite.db).
		Last(accessToken).RecordNotFound())

	// Check the response
	expected := &oauth.AccessTokenResponse{
		UserID:       accessToken.UserID.String,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    tokentypes.Bearer,
		Scope:        "read_write",
		RefreshToken: "test_token",
	}
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}
