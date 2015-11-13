package oauth

import (
	// "log"
	// "time"

	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

// func (suite *OauthTestSuite) TestGetOrCreateRefreshToken() {
// 	var refreshToken *RefreshToken
// 	var err error
//
// 	// Without user
// 	refreshToken, err = suite.service.getOrCreateRefreshToken(
// 		suite.client,
// 		nil,
// 		"foo bar",
// 	)
// 	assert.Nil(suite.T(), err)
//
// 	// With User
// 	refreshToken, err = suite.service.getOrCreateRefreshToken(
// 		suite.client,
// 		suite.user,
// 		"foo bar",
// 	)
// 	assert.Nil(suite.T(), err)
//
// 	log.Print(refreshToken)
// }

func (suite *OauthTestSuite) TestGetValidRefreshTokenNotFound() {
	refreshToken, err := s.getValidRefreshToken(
		"bogus",      // refresh token
		suite.client, // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Check the error
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Refresh token not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidRefreshTokenExpired() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	refreshToken, err := s.getValidRefreshToken(
		"test_refresh_token", // refresh token
		suite.client,         // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Check the error
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Refresh token expired", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidRefreshToken() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_refresh_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	refreshToken, err := s.getValidRefreshToken(
		"test_refresh_token", // refresh token
		suite.client,         // client
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Refresh token should be returned
	assert.NotNil(suite.T(), refreshToken)
}
