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

func (suite *OauthTestSuite) TestPasswordGrant() {
	// Make a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/oauth2/api/v1/tokens", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.PostForm = url.Values{
		"grant_type": {"password"},
		"username":   {"test_username"},
		"password":   {"test_password"},
	}

	w := httptest.NewRecorder()
	suite.service.passwordGrant(w, r, suite.client)

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
		"scope":         "foo bar",
		"refresh_token": refreshToken.Token,
	})
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
}
