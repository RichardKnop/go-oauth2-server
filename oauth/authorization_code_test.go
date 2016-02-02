package oauth

import (
	"database/sql/driver"
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGrantAuthorizationCode() {
	var (
		authorizationCode *AuthorizationCode
		err               error
		codes             []*AuthorizationCode
		v                 driver.Value
	)

	// Grant an authorization code
	authorizationCode, err = suite.service.GrantAuthorizationCode(
		suite.clients[0],              // client
		suite.users[0],                // user
		"redirect URI doesn't matter", // redirect URI
		"scope doesn't matter",        // scope
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
		assert.Equal(suite.T(), int64(suite.clients[0].ID), v)

		// User id should be set
		assert.True(suite.T(), codes[0].UserID.Valid)
		v, err = codes[0].UserID.Value()
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), int64(suite.users[0].ID), v)
	}
}

func (suite *OauthTestSuite) TestGetValidAuthorizationCode() {
	// Insert an expired test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:      "test_expired_code",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test authorization code
	if err := suite.db.Create(&AuthorizationCode{
		Code:      "test_code",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error; err != nil {
		log.Fatal(err)
	}

	var (
		authorizationCode *AuthorizationCode
		err               error
	)

	// Test passing an empty code
	authorizationCode, err = suite.service.getValidAuthorizationCode(
		"",               // authorization code
		suite.clients[0], // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAuthorizationCodeNotFound, err)
	}

	// Test passing a bogus code
	authorizationCode, err = suite.service.getValidAuthorizationCode(
		"bogus",          // authorization code
		suite.clients[0], // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAuthorizationCodeNotFound, err)
	}

	// Test passing an expired code
	authorizationCode, err = suite.service.getValidAuthorizationCode(
		"test_expired_code", // authorization code
		suite.clients[0],    // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errAuthorizationCodeExpired, err)
	}

	// Test passing a valid code
	authorizationCode, err = suite.service.getValidAuthorizationCode(
		"test_code",      // authorization code
		suite.clients[0], // client
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct authorization code object should be returned
	assert.NotNil(suite.T(), authorizationCode)
	assert.Equal(suite.T(), "test_code", authorizationCode.Code)
}
