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
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Form = url.Values{
		"grant_type": {"password"},
		"username":   {"test@user"},
		"password":   {"test_password"},
		"scope":      {"read_write"},
	}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.passwordGrant(w, r, suite.clients[0])

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the correct data was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.First(accessToken).RecordNotFound())
	refreshToken := new(RefreshToken)
	assert.False(suite.T(), suite.db.First(refreshToken).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&AccessTokenResponse{
		ID:           accessToken.ID,
		AccessToken:  accessToken.Token,
		ExpiresIn:    3600,
		TokenType:    TokenType,
		Scope:        "read_write",
		RefreshToken: refreshToken.Token,
	})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
}
