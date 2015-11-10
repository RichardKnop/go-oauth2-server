package oauth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestClientCredentialsGrant() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"bar qux"},
	}

	w := httptest.NewRecorder()
	suite.service.clientCredentialsGrant(w, r, suite.client)

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := AccessToken{}
	assert.False(suite.T(), suite.db.First(&accessToken).RecordNotFound())
	refreshToken := RefreshToken{}
	assert.False(suite.T(), suite.db.First(&refreshToken).RecordNotFound())

	// Check the response body
	expected, _ := json.Marshal(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    3600,
		"token_type":    "Bearer",
		"scope":         "bar qux",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
}
