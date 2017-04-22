package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/RichardKnop/go-oauth2-server/models"
	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/oauth/tokentypes"
	testutil "github.com/RichardKnop/go-oauth2-server/test-util"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/RichardKnop/uuid"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestNewIntrospectResponseFromAccessToken() {
	MG := models.MyGormModel{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}

	accessToken := &models.OauthAccessToken{
		MyGormModel: MG,
		Token:       "test_token_introspect_1",
		ExpiresAt:   time.Now().UTC().Add(+10 * time.Second),
		ClientID:    util.StringOrNull(string(suite.clients[0].ID)),
		UserID:      util.StringOrNull(string(suite.users[0].ID)),
		Scope:       "read_write",
	}
	expected := &oauth.IntrospectResponse{
		Active:    true,
		Scope:     accessToken.Scope,
		TokenType: tokentypes.Bearer,
		ExpiresAt: int(accessToken.ExpiresAt.Unix()),
		ClientID:  suite.clients[0].Key,
		Username:  suite.users[0].Username,
	}

	actual, err := suite.service.NewIntrospectResponseFromAccessToken(accessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)

	accessToken.ClientID = util.StringOrNull("")
	expected.ClientID = ""
	actual, err = suite.service.NewIntrospectResponseFromAccessToken(accessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)

	accessToken.UserID = util.StringOrNull("")
	expected.Username = ""
	actual, err = suite.service.NewIntrospectResponseFromAccessToken(accessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)
}

func (suite *OauthTestSuite) TestNewIntrospectResponseFromRefreshToken() {
	MG := models.MyGormModel{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
	}

	refreshToken := &models.OauthRefreshToken{
		MyGormModel: MG,
		Token:       "test_token_introspect_1",
		ExpiresAt:   time.Now().UTC().Add(+10 * time.Second),
		ClientID:    util.StringOrNull(string(suite.clients[0].ID)),
		UserID:      util.StringOrNull(string(suite.users[0].ID)),
		Scope:       "read_write",
	}
	expected := &oauth.IntrospectResponse{
		Active:    true,
		Scope:     refreshToken.Scope,
		TokenType: tokentypes.Bearer,
		ExpiresAt: int(refreshToken.ExpiresAt.Unix()),
		ClientID:  suite.clients[0].Key,
		Username:  suite.users[0].Username,
	}

	actual, err := suite.service.NewIntrospectResponseFromRefreshToken(refreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)

	refreshToken.ClientID = util.StringOrNull("")
	expected.ClientID = ""
	actual, err = suite.service.NewIntrospectResponseFromRefreshToken(refreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)

	refreshToken.UserID = util.StringOrNull("")
	expected.Username = ""
	actual, err = suite.service.NewIntrospectResponseFromRefreshToken(refreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, actual)
}

func (suite *OauthTestSuite) TestHandleIntrospectMissingToken() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrTokenMissing.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectInvailidTokenHint() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{"token": {"token"}, "token_type_hint": {"wrong"}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrTokenHintInvalid.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectAccessToken() {
	// Insert a test access token with a user
	accessToken := &models.OauthAccessToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(accessToken).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{
		"token":           {accessToken.Token},
		"token_type_hint": {oauth.AccessTokenHint},
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	expected, err := suite.service.NewIntrospectResponseFromAccessToken(accessToken)
	assert.NoError(suite.T(), err)
	testutil.TestResponseObject(suite.T(), w, expected, 200)

	// With incorrect token hint
	r.PostForm = url.Values{
		"token":           {accessToken.Token},
		"token_type_hint": {oauth.RefreshTokenHint},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRefreshTokenNotFound.Error(),
		404,
	)

	// Without token hint
	r.PostForm = url.Values{
		"token": {accessToken.Token},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	expected, err = suite.service.NewIntrospectResponseFromAccessToken(accessToken)
	assert.NoError(suite.T(), err)
	testutil.TestResponseObject(suite.T(), w, expected, 200)
}

func (suite *OauthTestSuite) TestHandleIntrospectRefreshToken() {
	// Insert a test refresh token with a user
	refreshToken := &models.OauthRefreshToken{
		MyGormModel: models.MyGormModel{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().UTC().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(refreshToken).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{
		"token":           {refreshToken.Token},
		"token_type_hint": {oauth.RefreshTokenHint},
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	expected, err := suite.service.NewIntrospectResponseFromRefreshToken(refreshToken)
	assert.NoError(suite.T(), err)
	testutil.TestResponseObject(suite.T(), w, expected, 200)

	// With incorrect token hint
	r.PostForm = url.Values{
		"token":           {refreshToken.Token},
		"token_type_hint": {oauth.AccessTokenHint},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAccessTokenNotFound.Error(),
		404,
	)

	// Without token hint
	r.PostForm = url.Values{
		"token": {refreshToken.Token},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAccessTokenNotFound.Error(),
		404,
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectInactiveToken() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With access token hint
	r.PostForm = url.Values{
		"token":           {"unexisting_token"},
		"token_type_hint": {oauth.AccessTokenHint},
	}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAccessTokenNotFound.Error(),
		404,
	)

	// With refresh token hint
	r.PostForm = url.Values{
		"token":           {"unexisting_token"},
		"token_type_hint": {oauth.RefreshTokenHint},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrRefreshTokenNotFound.Error(),
		404,
	)

	// Without token hint
	r.PostForm = url.Values{
		"token": {"unexisting_token"},
	}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrAccessTokenNotFound.Error(),
		404,
	)
}
