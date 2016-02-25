package oauth

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticate() {
	var (
		accessToken *AccessToken
		err         error
	)

	// Insert some test access tokens
	testAccessTokens := []*AccessToken{
		// Expired access token
		&AccessToken{
			Token:     "test_expired_token",
			ExpiresAt: time.Now().Add(-10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		// Access token without a user
		&AccessToken{
			Token:     "test_client_token",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
		},
		// Access token with a user
		&AccessToken{
			Token:     "test_user_token",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
	}
	for _, testAccessToken := range testAccessTokens {
		err := suite.db.Create(testAccessToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

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
