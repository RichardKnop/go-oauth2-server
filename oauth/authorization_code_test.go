package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetValidAuthorizationCodeNotFound() {
	authorizationCode, err := suite.service.getValidAuthorizationCode(
		"bogus",      // authorization code
		suite.client, // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Authorization code not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidAuthorizationCodeExpired() {
	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:      "test_authorization_code",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	authorizationCode, err := suite.service.getValidAuthorizationCode(
		"test_authorization_code", // authorization code
		suite.client,              // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Authorization code expired", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidAuthorizationCode() {
	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:      "test_code",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	authorizationCode, err := suite.service.getValidAuthorizationCode(
		"test_code",  // authorization code
		suite.client, // client
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct authorization code object should be returned
	assert.NotNil(suite.T(), authorizationCode)
	assert.Equal(suite.T(), "test_code", authorizationCode.Code)
}
