package oauth

import (
	"database/sql/driver"
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGrantAuthorizationCode() {
	var authorizationCode *AuthorizationCode
	var err error
	var codes []*AuthorizationCode
	var v driver.Value

	// Grant an authorization code
	authorizationCode, err = suite.service.GrantAuthorizationCode(
		suite.client,
		suite.user,
		"redirect URI doesn't matter",
		"scope doesn't matter",
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct authorization code object should be returned
	if assert.NotNil(suite.T(), authorizationCode) {
		// Fetch all access tokens
		suite.service.db.Preload("Client").Preload("User").Find(&codes)

		// There should be just one right now
		assert.Equal(suite.T(), 1, len(codes))

		// And the code should match the one returned by the grant method
		assert.Equal(suite.T(), codes[0].Code, authorizationCode.Code)

		// Client id should be set
		assert.True(suite.T(), codes[0].ClientID.Valid)
		v, err = codes[0].ClientID.Value()
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), int64(suite.client.ID), v)

		// User id should be set
		assert.True(suite.T(), codes[0].UserID.Valid)
		v, err = codes[0].UserID.Value()
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), int64(suite.user.ID), v)
	}
}

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
		Code:      "test_code",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
	}).Error; err != nil {
		log.Fatal(err)
	}

	authorizationCode, err := suite.service.getValidAuthorizationCode(
		"test_code",  // authorization code
		suite.client, // client
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
