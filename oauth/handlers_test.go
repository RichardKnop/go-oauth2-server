package oauth

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestHandleTokensClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{"grant_type": {"client_credentials"}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.tokensHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", errClientAuthenticationRequired.Error()),
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

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.tokensHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 400, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", errInvalidGrantType.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}

func (suite *OauthTestSuite) TestHandleIntrospectClientAuthenticationRequired() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{"token": {"token"}}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.introspectHandler(w, r)

	// Check the status code
	assert.Equal(suite.T(), 401, w.Code)

	// Check the response body
	assert.Equal(
		suite.T(),
		fmt.Sprintf("{\"error\":\"%s\"}", errClientAuthenticationRequired.Error()),
		strings.TrimSpace(w.Body.String()),
	)
}
