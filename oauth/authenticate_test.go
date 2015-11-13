package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticateNotFound() {
	accessToken, err := s.Authenticate("bogus")

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
		Scope:     "foo",
	}).Error; err != nil {
		log.Fatal(err)
	}

	accessToken, err := s.Authenticate("test_token")

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
		Scope:     "foo",
	}).Error; err != nil {
		log.Fatal(err)
	}

	accessToken, err := s.Authenticate("test_token")

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct access token object should be returned
	if assert.NotNil(suite.T(), accessToken) {
		assert.Equal(suite.T(), "test_token", accessToken.Token)
	}
}
