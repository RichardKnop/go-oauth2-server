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

func (suite *OauthTestSuite) TestIntrospectResponseAccessToken() {
	at := &oauth.AccessToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		ClientID:  util.PositiveIntOrNull(int64(suite.clients[0].ID)),
		UserID:    util.PositiveIntOrNull(int64(suite.users[0].ID)),
		Scope:     "read_write",
	}
	ir := &oauth.IntrospectResponse{
		Active:    true,
		Scope:     at.Scope,
		TokenType: oauth.TokenType,
		ExpiresAt: int(at.ExpiresAt.Unix()),
		ClientID:  suite.clients[0].Key,
		Username:  suite.users[0].Username,
	}
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))

	at.ClientID = util.PositiveIntOrNull(0)
	ir.ClientID = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))

	at.UserID = util.PositiveIntOrNull(0)
	ir.Username = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))
}

func (suite *OauthTestSuite) TestIntrospectResponseRefreshToken() {
	rt := &oauth.RefreshToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		ClientID:  util.PositiveIntOrNull(int64(suite.clients[0].ID)),
		UserID:    util.PositiveIntOrNull(int64(suite.users[0].ID)),
		Scope:     "read_write",
	}
	ir := &oauth.IntrospectResponse{
		Active:    true,
		Scope:     rt.Scope,
		TokenType: oauth.TokenType,
		ExpiresAt: int(rt.ExpiresAt.Unix()),
		ClientID:  suite.clients[0].Key,
		Username:  suite.users[0].Username,
	}
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))

	rt.ClientID = util.PositiveIntOrNull(0)
	ir.ClientID = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))

	rt.UserID = util.PositiveIntOrNull(0)
	ir.Username = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))
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

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrTokenMissing.Error()),
		strings.TrimSpace(w.Body.String()),
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

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrTokenHintInvalid.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectAccessToken() {
	// Insert a test access token with a user
	at := &oauth.AccessToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(at).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{"token": {at.Token}, "token_type_hint": {oauth.AccessTokenHint}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(suite.service.IntrospectResponseAccessToken(at))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With incorrect token hint
	r.PostForm = url.Values{"token": {at.Token}, "token_type_hint": {oauth.RefreshTokenHint}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseAccessToken(at))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {at.Token}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseAccessToken(at))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *OauthTestSuite) TestHandleIntrospectRefreshToken() {
	// Insert a test access token with a user
	rt := &oauth.RefreshToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(rt).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{"token": {rt.Token}, "token_type_hint": {oauth.RefreshTokenHint}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(suite.service.IntrospectResponseRefreshToken(rt))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With incorrect token hint
	r.PostForm = url.Values{"token": {rt.Token}, "token_type_hint": {oauth.AccessTokenHint}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseRefreshToken(rt))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {rt.Token}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseRefreshToken(rt))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}

func (suite *OauthTestSuite) TestHandleIntrospectInactiveToken() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With access token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}, "token_type_hint": {oauth.AccessTokenHint}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(new(oauth.IntrospectResponse))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With refresh token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}, "token_type_hint": {oauth.RefreshTokenHint}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(new(oauth.IntrospectResponse))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}}

	// Serve the request
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(new(oauth.IntrospectResponse))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
