package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticateNotFound() {
	accessToken, err := suite.service.Authenticate("bogus")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Access token not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthenticateExpired() {
	// Insert a test access token
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
	}).Error; err != nil {
		log.Fatal(err)
	}

	accessToken, err := suite.service.Authenticate("test_token")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Access token expired", err.Error())
	}
}

func (suite *OauthTestSuite) TestAuthenticate() {
	// Insert a test access token
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
	}).Error; err != nil {
		log.Fatal(err)
	}

	accessToken, err := suite.service.Authenticate("test_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_token", accessToken.Token)
		assert.Equal(suite.T(), "test_client", accessToken.Client.ClientID)
		assert.Equal(suite.T(), "test@username", accessToken.User.Username)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)
}
