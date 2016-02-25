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

func (suite *OauthTestSuite) TestIntrospectResponseAccessToken() {
	at := &AccessToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	ir := &IntrospectResponse{
		Active:    true,
		Scope:     at.Scope,
		TokenType: TokenType,
		ExpiresAt: int(at.ExpiresAt.Unix()),
		ClientID:  at.Client.Key,
		Username:  at.User.Username,
	}
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))

	at.Client = nil
	ir.ClientID = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))

	at.User = nil
	ir.Username = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseAccessToken(at))
}

func (suite *OauthTestSuite) TestIntrospectResponseRefreshToken() {
	rt := &RefreshToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	ir := &IntrospectResponse{
		Active:    true,
		Scope:     rt.Scope,
		TokenType: TokenType,
		ExpiresAt: int(rt.ExpiresAt.Unix()),
		ClientID:  rt.Client.Key,
		Username:  rt.User.Username,
	}
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))

	rt.Client = nil
	ir.ClientID = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))

	rt.User = nil
	ir.Username = ""
	assert.Equal(suite.T(), ir, suite.service.IntrospectResponseRefreshToken(rt))
}

func (suite *OauthTestSuite) TestHandleIntrospectMissingToken() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", errTokenMissing.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectInvailidTokenHint() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")
	r.PostForm = url.Values{"token": {"token"}, "token_type_hint": {"wrong"}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", errTokenHintInvalid.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectAccessToken() {
	// Insert a test access token with a user
	at := &AccessToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(at).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{"token": {at.Token}, "token_type_hint": {accessTokenHint}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(suite.service.IntrospectResponseAccessToken(at))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With incorrect token hint
	r.PostForm = url.Values{"token": {at.Token}, "token_type_hint": {refreshTokenHint}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseAccessToken(at))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {at.Token}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

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
	rt := &RefreshToken{
		Token:     "test_token_introspect_1",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
		Scope:     "read_write",
	}
	err := suite.db.Create(rt).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With correct token hint
	r.PostForm = url.Values{"token": {rt.Token}, "token_type_hint": {refreshTokenHint}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(suite.service.IntrospectResponseRefreshToken(rt))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With incorrect token hint
	r.PostForm = url.Values{"token": {rt.Token}, "token_type_hint": {accessTokenHint}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(suite.service.IntrospectResponseRefreshToken(rt))
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {rt.Token}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

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
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client_1", "test_secret")

	// With access token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}, "token_type_hint": {accessTokenHint}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err := json.Marshal(&IntrospectResponse{})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// With refresh token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}, "token_type_hint": {refreshTokenHint}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(&IntrospectResponse{})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}

	// Without token hint
	r.PostForm = url.Values{"token": {"unexisting_token"}}

	// And run the function we want to test
	w = httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)
	// Check the response body
	expected, err = json.Marshal(&IntrospectResponse{})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
