package oauth

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGrantAuthorizationCode() {
	var (
		authorizationCode *AuthorizationCode
		err               error
		codes             []*AuthorizationCode
	)

	// Grant an authorization code
	authorizationCode, err = suite.service.GrantAuthorizationCode(
		suite.clients[0], // client
		suite.users[0],   // user
		3600,             // expires in
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

		// Client ID should be set
		assert.True(suite.T(), codes[0].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), codes[0].ClientID.Int64)

		// User ID should be set
		assert.True(suite.T(), codes[0].UserID.Valid)
		assert.Equal(suite.T(), int64(suite.users[0].ID), codes[0].UserID.Int64)
	}
}

func (suite *OauthTestSuite) TestGetValidAuthorizationCode() {
	var (
		authorizationCode *AuthorizationCode
		err               error
	)

	// Insert some test authorization codes
	testAuthorizationCodes := []*AuthorizationCode{
		// Expired authorization code
		&AuthorizationCode{
			Code:      "test_expired_code",
			ExpiresAt: time.Now().Add(-10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		// Authorization code
		&AuthorizationCode{
			Code:      "test_code",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
	}
	for _, testAuthorizationCode := range testAuthorizationCodes {
		err := suite.db.Create(testAuthorizationCode).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Test passing an empty code
	authorizationCode, err = suite.service.getValidAuthorizationCode(
		"",               // authorization code
		suite.clients[0], // client
	)

	// Authorization code should be nil
	assert.Nil(suite.T(), authorizationCode)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), ErrAuthorizationCodeNotFound, err)
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
		assert.Equal(suite.T(), ErrAuthorizationCodeNotFound, err)
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
		assert.Equal(suite.T(), ErrAuthorizationCodeExpired, err)
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
