package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestClientCredentialsGrant() {
	// Prepare a request
	r, err := http.NewRequest("POST", "http://1.2.3.4/something", nil)
	assert.NoError(suite.T(), err, "Request setup should not get an error")
	r.Form = url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {"read_write"},
	}

	// And run the function we want to test
	w := httptest.NewRecorder()
	suite.service.clientCredentialsGrant(w, r, suite.clients[0])

	// Check the status code
	assert.Equal(suite.T(), 200, w.Code)

	// Check the access token was inserted
	accessToken := new(AccessToken)
	assert.False(suite.T(), suite.db.Preload("Client").Preload("User").
		First(accessToken).RecordNotFound())

	// Client credentials grant does not produce refresh token
	assert.True(suite.T(), suite.db.Preload("Client").Preload("User").
		First(new(RefreshToken)).RecordNotFound())

	// Check the response body
	expected, err := json.Marshal(&AccessTokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   3600,
		TokenType:   TokenType,
		Scope:       "read_write",
	})
	if assert.NoError(suite.T(), err, "JSON marshalling failed") {
		assert.Equal(suite.T(), string(expected), strings.TrimSpace(w.Body.String()))
	}
}
