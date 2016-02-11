package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticate() {
	// Insert an expired test access token
	if err := suite.db.Create(&AccessToken{
		Token:     "test_expired_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test client access token
	if err := suite.db.Create(&AccessToken{
		Token:     "test_client_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test user access token
	if err := suite.db.Create(&AccessToken{
		Token:     "test_user_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error; err != nil {
		log.Fatal(err)
	}

	var (
		accessToken *AccessToken
		err         error
	)

	// Test passing an empty token
	accessToken, err = suite.service.Authenticate("")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAccessTokenNotFound, err)
	}

	// Test passing a bogus token
	accessToken, err = suite.service.Authenticate("bogus")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAccessTokenNotFound, err)
	}

	// Test passing an expired token
	accessToken, err = suite.service.Authenticate("test_expired_token")

	// Access token should be nil
	assert.Nil(suite.T(), accessToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAccessTokenExpired, err)
	}

	// Test passing a valid client token
	accessToken, err = suite.service.Authenticate("test_client_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_client_token", accessToken.Token)
		assert.Equal(suite.T(), "test_client_1", accessToken.Client.Key)
		assert.Nil(suite.T(), accessToken.User)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Test passing a valid user token
	accessToken, err = suite.service.Authenticate("test_user_token")

	// Correct access token should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_user_token", accessToken.Token)
		assert.Equal(suite.T(), "test_client_1", accessToken.Client.Key)
		assert.Equal(suite.T(), "test@superuser", accessToken.User.Username)
	}

	// Error should be nil
	assert.Nil(suite.T(), err)
}
