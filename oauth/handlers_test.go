package oauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/RichardKnop/go-oauth2-server/oauth"
	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestTokensHandlerClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{"grant_type": {"client_credentials"}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrInvalidClientIDOrSecret.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestTokensHandlerInvalidGrantType() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/tokens", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{"grant_type": {"bogus"}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrInvalidGrantType.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestIntrospectHandlerClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/v1/oauth/introspect", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.PostForm = url.Values{"token": {"token"}}

	// And serve the request
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, r)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", oauth.ErrInvalidClientIDOrSecret.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}
