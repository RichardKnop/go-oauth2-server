package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestAuthenticateNotFound() {
	err := suite.service.Authenticate("bogus")

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

	err := suite.service.Authenticate("test_token")

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

	// Error should be nil
	assert.Nil(suite.T(), suite.service.Authenticate("test_token"))
}
