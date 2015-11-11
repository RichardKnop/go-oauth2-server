package oauth

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestHandleTokensClientAuthenticationRequired() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{"grant_type": {"client_credentials"}}

	w := httptest.NewRecorder()
	handleTokens(w, r)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Client authentication required\"}",
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestHandleTokensInvalidGrantType() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.SetBasicAuth("test_client", "test_secret")
	r.PostForm = url.Values{"grant_type": {"bogus"}}

	w := httptest.NewRecorder()
	handleTokens(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(), "{\"error\":\"Invalid grant type\"}",
		strings.TrimSpace(w.Body.String()),
	)
}
