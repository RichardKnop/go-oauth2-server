package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenCreatesNew() {
	// Since there is no token, a new one should be created and returned
	refreshToken, err := suite.service.GetOrCreateRefreshToken(
		suite.client,
		suite.user,
		"foo",
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// There should be just one refresh token
	var count int
	var tokens []*RefreshToken
	s.db.Where(RefreshToken{
		ClientID: clientIDOrNull(suite.client),
		UserID:   userIDOrNull(suite.user),
	}).Find(&tokens).Count(&count)
	assert.Equal(suite.T(), 1, count)

	// Correct refresh token object should be returned
	assert.NotNil(suite.T(), refreshToken)
	assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)
}

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenReturnsExisting() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Since the current token is valid, this should just return it
	refreshToken, err := suite.service.GetOrCreateRefreshToken(
		suite.client,
		suite.user,
		"foo",
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// The valid token should NOT have been deleted
	found := !s.db.Where(RefreshToken{
		Token: "test_token",
	}).First(new(RefreshToken)).RecordNotFound()
	assert.True(suite.T(), found)

	// There should be just one refresh token in the database
	var count int
	var tokens []*RefreshToken
	s.db.Where(RefreshToken{
		ClientID: clientIDOrNull(suite.client),
		UserID:   userIDOrNull(suite.user),
	}).Find(&tokens).Count(&count)
	assert.Equal(suite.T(), 1, count)

	// Correct refresh token object should be returned
	assert.NotNil(suite.T(), refreshToken)
	assert.Equal(suite.T(), "test_token", refreshToken.Token)
	assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)
}

func (suite *OauthTestSuite) TestGetOrCreateRefreshTokenDeletesExpired() {
	// Insert an expired test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Since the current token is expired, this should delete it
	// and create and return a new one
	refreshToken, err := suite.service.GetOrCreateRefreshToken(
		suite.client,
		suite.user,
		"foo",
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// The expired token should have been deleted
	notFound := s.db.Where(RefreshToken{
		Token: "test_token",
	}).First(new(RefreshToken)).RecordNotFound()
	assert.True(suite.T(), notFound)

	// There should be just one refresh token
	var count int
	var tokens []*RefreshToken
	s.db.Where(RefreshToken{
		ClientID: clientIDOrNull(suite.client),
		UserID:   userIDOrNull(suite.user),
	}).Find(&tokens).Count(&count)
	assert.Equal(suite.T(), 1, count)

	// Correct refresh token object should be returned
	assert.NotNil(suite.T(), refreshToken)
	assert.NotEqual(suite.T(), "test_token", refreshToken.Token)
	assert.Equal(suite.T(), tokens[0].Token, refreshToken.Token)
}

func (suite *OauthTestSuite) TestGetValidRefreshTokenNotFound() {
	refreshToken, err := s.getValidRefreshToken(
		"bogus",      // refresh token
		suite.client, // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Refresh token not found", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidRefreshTokenExpired() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	refreshToken, err := s.getValidRefreshToken(
		"test_token", // refresh token
		suite.client, // client
	)

	// Refresh token should be nil
	assert.Nil(suite.T(), refreshToken)

	// Correct error should be returned
	if assert.NotNil(suite.T(), err) {
		assert.Equal(suite.T(), "Refresh token expired", err.Error())
	}
}

func (suite *OauthTestSuite) TestGetValidRefreshToken() {
	// Insert a test refresh token
	if err := suite.db.Create(&RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	refreshToken, err := s.getValidRefreshToken(
		"test_token", // refresh token
		suite.client, // client
	)

	// Error should be nil
	assert.Nil(suite.T(), err)

	// Correct refresh token object should be returned
	assert.NotNil(suite.T(), refreshToken)
	assert.Equal(suite.T(), "test_token", refreshToken.Token)
}
