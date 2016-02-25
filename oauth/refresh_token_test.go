package oauth

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenCreatesNew() {
	var (
		refreshToken *RefreshToken
		err          error
		tokens       []*RefreshToken
	)

	// Since there is no client only token,
	// a new one should be created and returned
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		nil,              // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be nil
	if assert.Nil(suite.T(), err) {
		// Fetch all refresh tokens
		suite.service.db.Preload("Client").Preload("User").Find(&tokens)

		// There should be just one right now
		assert.Equal(suite.T(), 1, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[0].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[0].ClientID.Int64)

		// User ID should be nil
		assert.False(suite.T(), tokens[0].UserID.Valid)
	}

	// Since there is no user specific token,
	// a new one should be created and returned
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		suite.users[0],   // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be nil
	if assert.Nil(suite.T(), err) {
		// Fetch all refresh tokens
		suite.service.db.Preload("Client").Preload("User").Find(&tokens)

		// There should be 2 tokens now
		assert.Equal(suite.T(), 2, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[1].Token, refreshToken.Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[1].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[1].ClientID.Int64)

		// User ID should be set
		assert.True(suite.T(), tokens[1].UserID.Valid)
		assert.Equal(suite.T(), int64(suite.users[0].ID), tokens[1].UserID.Int64)
	}
}

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenReturnsExisting() {
	var (
		refreshToken *RefreshToken
		err          error
		tokens       []*RefreshToken
	)

	// Insert an access token without a user
	err = suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Since the current client only token is valid, this should just return it
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		nil,              // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct refresh token should be returned
	if assert.NotNil(suite.T(), refreshToken) {
		// Fetch all refresh tokens
		suite.service.db.Preload("Client").Preload("User").Find(&tokens)

		// There should be just one token right now
		assert.Equal(suite.T(), 1, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)
		assert.Equal(suite.T(), "test_token", refreshToken.Token)
		assert.Equal(suite.T(), "test_token", tokens[0].Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[0].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[0].ClientID.Int64)

		// User ID should be nil
		assert.False(suite.T(), tokens[0].UserID.Valid)
	}

	// Insert an access token with a user
	err = suite.db.Create(&RefreshToken{
		Token:     "test_token2",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Since the current user specific only token is valid,
	// this should just return it
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		suite.users[0],   // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be Nil
	assert.Nil(suite.T(), err)

	// Correct refresh token should be returned
	if assert.NotNil(suite.T(), refreshToken) {
		// Fetch all refresh tokens
		suite.service.db.Preload("Client").Preload("User").Find(&tokens)

		// There should be 2 tokens now
		assert.Equal(suite.T(), 2, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[1].Token, refreshToken.Token)
		assert.Equal(suite.T(), "test_token2", refreshToken.Token)
		assert.Equal(suite.T(), "test_token2", tokens[1].Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[1].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[1].ClientID.Int64)

		// User ID should be set
		assert.True(suite.T(), tokens[1].UserID.Valid)
		assert.Equal(suite.T(), int64(suite.users[0].ID), tokens[1].UserID.Int64)
	}
}

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenDeletesExpired() {
	var (
		refreshToken *RefreshToken
		err          error
		tokens       []*RefreshToken
	)

	// Insert an expired client only test refresh token
	err = suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.clients[0],
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Since the current client only token is expired,
	// this should delete it and create and return a new one
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		nil,              // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be nil
	if assert.Nil(suite.T(), err) {
		// Fetch all refresh tokens
		suite.service.db.Unscoped().Preload("Client").Preload("User").Find(&tokens)

		// There should be just one token right now
		assert.Equal(suite.T(), 1, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)
		assert.NotEqual(suite.T(), "test_token", refreshToken.Token)
		assert.NotEqual(suite.T(), "test_token", tokens[0].Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[0].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[0].ClientID.Int64)

		// User ID should be nil
		assert.False(suite.T(), tokens[0].UserID.Valid)
	}

	// Insert an expired user specific test refresh token
	err = suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.clients[0],
		User:      suite.users[0],
	}).Error
	assert.NoError(suite.T(), err, "Inserting test data failed")

	// Since the current user specific token is expired,
	// this should delete it and create and return a new one
	refreshToken, err = suite.service.GetOrCreateRefreshToken(
		suite.clients[0], // client
		suite.users[0],   // user
		3600,             // expires in
		"read_write",     // scope
	)

	// Error should be nil
	if assert.Nil(suite.T(), err) {
		// Fetch all refresh tokens
		suite.service.db.Unscoped().Preload("Client").Preload("User").Find(&tokens)

		// There should be 2 tokens now
		assert.Equal(suite.T(), 2, len(tokens))

		// Correct refresh token object should be returned
		assert.NotNil(suite.T(), refreshToken)
		assert.Equal(suite.T(), tokens[1].Token, refreshToken.Token)
		assert.NotEqual(suite.T(), "test_token", refreshToken.Token)
		assert.NotEqual(suite.T(), "test_token", tokens[1].Token)

		// Client ID should be set
		assert.True(suite.T(), tokens[1].ClientID.Valid)
		assert.Equal(suite.T(), int64(suite.clients[0].ID), tokens[1].ClientID.Int64)

		// User ID should be set
		assert.True(suite.T(), tokens[1].UserID.Valid)
		assert.Equal(suite.T(), int64(suite.users[0].ID), tokens[1].UserID.Int64)
	}
}

func (suite *OauthTestSuite) TestGetValidRefreshToken() {
	var (
		refreshToken *RefreshToken
		err          error
	)

	// Insert some test refresh tokens
	testRefreshTokens := []*RefreshToken{
		// Expired test refresh token
		&RefreshToken{
			Token:     "test_expired_token",
			ExpiresAt: time.Now().Add(-10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
		// Refresh token
		&RefreshToken{
			Token:     "test_token",
			ExpiresAt: time.Now().Add(+10 * time.Second),
			Client:    suite.clients[0],
			User:      suite.users[0],
		},
	}
	for _, testRefreshToken := range testRefreshTokens {
		err := suite.db.Create(testRefreshToken).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Test passing an empty token
	refreshToken, err = suite.service.GetValidRefreshToken(
		"",               // refresh token
		suite.clients[0], // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errRefreshTokenNotFound, err)
	}

	// Test passing a bogus token
	refreshToken, err = suite.service.GetValidRefreshToken(
		"bogus",          // refresh token
		suite.clients[0], // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errRefreshTokenNotFound, err)
	}

	// Test passing an expired token
	refreshToken, err = suite.service.GetValidRefreshToken(
		"test_expired_token", // refresh token
		suite.clients[0],     // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), errRefreshTokenExpired, err)
	}

	// Test passing a valid token
	refreshToken, err = suite.service.GetValidRefreshToken(
		"test_token",     // refresh token
		suite.clients[0], // client
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct refresh token object should be returned
	assert.NotNil(suite.T(), refreshToken)
	assert.Equal(suite.T(), "test_token", refreshToken.Token)
}
