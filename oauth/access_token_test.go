package oauth

import (
	"log"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *OauthTestSuite) TestDeleteExpiredAccessTokens() {
	// Insert an expired test access token with a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_1",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert an expired test access token without a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_2",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    suite.client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test access token with a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_3",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		User:      suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test access token without a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_4",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    suite.client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	var notFound bool
	var existingTokens []string

	// This should only delete test_token_1
	suite.service.deleteExpiredAccessTokens(suite.client, suite.user)

	// Check the test_token_1 was deleted
	notFound = suite.db.Where(AccessToken{Token: "test_token_1"}).
		First(&AccessToken{}).RecordNotFound()
	assert.True(suite.T(), notFound)

	// Check the other three tokens are still around
	existingTokens = []string{
		"test_token_2",
		"test_token_3",
		"test_token_4",
	}
	for _, token := range existingTokens {
		notFound = suite.db.Where(AccessToken{Token: token}).
			First(new(AccessToken)).RecordNotFound()
		assert.False(suite.T(), notFound)
	}

	// This should only delete test_token_2
	suite.service.deleteExpiredAccessTokens(suite.client, nil)

	// Check the test_token_2 was deleted
	notFound = suite.db.Where(AccessToken{Token: "test_token_2"}).
		First(new(AccessToken)).RecordNotFound()
	assert.True(suite.T(), notFound)

	// Check that last two tokens are still around
	existingTokens = []string{
		"test_token_3",
		"test_token_4",
	}
	for _, token := range existingTokens {
		notFound := suite.db.Where(AccessToken{Token: token}).
			First(new(AccessToken)).RecordNotFound()
		assert.False(suite.T(), notFound)
	}
}
