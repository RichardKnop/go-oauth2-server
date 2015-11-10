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
		Client:    *suite.client,
		User:      *suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert an expired test access token without a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_2",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test access token with a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_3",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.client,
		User:      *suite.user,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert a test access token without a user
	if err := suite.db.Create(&AccessToken{
		Token:     "test_token_4",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// This should only delete test_token_1
	suite.service.deleteExpiredAccessTokens(suite.client, suite.user)

	// Check the test_token_1 was deleted
	assert.True(
		suite.T(),
		suite.db.Where(&AccessToken{Token: "test_token_1"}).
			First(&AccessToken{}).RecordNotFound(),
	)

	// Check the other three tokens are still around
	for _, token := range []string{"test_token_2", "test_token_3", "test_token_4"} {
		assert.False(
			suite.T(),
			suite.db.Where(&AccessToken{Token: token}).
				First(&AccessToken{}).RecordNotFound(),
		)
	}

	// This should only delete test_token_2
	suite.service.deleteExpiredAccessTokens(suite.client, nil)

	// Check the test_token_2 was deleted
	assert.True(
		suite.T(),
		suite.db.Where(&AccessToken{Token: "test_token_2"}).
			First(&AccessToken{}).RecordNotFound(),
	)

	// Check that last two tokens are still around
	for _, token := range []string{"test_token_3", "test_token_4"} {
		assert.False(
			suite.T(),
			suite.db.Where(&AccessToken{Token: token}).
				First(&AccessToken{}).RecordNotFound(),
		)
	}
}

// func (suite *OauthTestSuite) TestGetOrCreateRefreshToken() {
// 	client := Client{ClientID: "test_client"}
// 	user := User{Username: "test_username"}
//
// 	var refreshToken *RefreshToken
// 	var err error
//
// 	// client_id = 1 AND user_id IS NULL
// 	refreshToken, err = getOrCreateRefreshToken(suite.DB, &client, nil, 1209600, "foo bar")
// 	assert.Nil(suite.T(), err)
//
// 	// client_id = 1 AND user_id = 2
// 	refreshToken, err = getOrCreateRefreshToken(suite.DB, &client, &user, 1209600, "foo bar")
// 	assert.Nil(suite.T(), err)
//
// 	log.Print(refreshToken)
// }
