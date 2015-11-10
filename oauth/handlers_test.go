package oauth

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestInvalidGrantType() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
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
