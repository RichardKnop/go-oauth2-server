package oauth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/RichardKnop/go-oauth2-server/test-util"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestTokensHandlerClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{"grant_type": {"client_credentials"}}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrInvalidClientIDOrSecret.Error(),
		401,
	)
}

func (suite *OauthTestSuite) TestTokensHandlerInvalidGrantType() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{"grant_type": {"bogus"}}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrInvalidGrantType.Error(),
		400,
	)
}

func (suite *OauthTestSuite) TestIntrospectHandlerClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{"token": {"token"}}

	// Serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the response
	testutil.TestResponseForError(
		suite.T(),
		w,
		oauth.ErrInvalidClientIDOrSecret.Error(),
		401,
	)
}
