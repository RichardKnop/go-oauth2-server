package oauth

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	client := Client{ClientID: "test_client"}
	user := User{Username: "test_username"}

	var accessToken *AccessToken

	// With user
	accessToken = newAccessToken(3600, &client, &user, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, "test_username", accessToken.User.Username,
		"Access token should belong to test_username",
	)

	// Without user
	accessToken = newAccessToken(3600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, uint(0), accessToken.User.ID,
		"Access token should not belong to a user",
	)
}

func TestNewRefreshToken(t *testing.T) {
	client := Client{ClientID: "test_client"}
	user := User{Username: "test_username"}

	var refreshToken *RefreshToken

	// With user
	refreshToken = newRefreshToken(3600, &client, &user, "doesn't matter")
	assert.Equal(
		t, "test_client", refreshToken.Client.ClientID,
		"Refresh token should belong to test_client",
	)
	assert.Equal(
		t, "test_username", refreshToken.User.Username,
		"Refresh token should belong to test_username",
	)

	// Without user
	refreshToken = newRefreshToken(3600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", refreshToken.Client.ClientID,
		"Refresh token should belong to test_client",
	)
	assert.Equal(
		t, uint(0), refreshToken.User.ID,
		"Refresh token should not belong to a user",
	)
}

func (suite *OauthTestSuite) TestDeleteExpiredAccessTokens() {
	// Insert expired test access token with user
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_1",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert expired test access token without user
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_2",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.Client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test access token with user
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_3",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test access token without user
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_4",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.Client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// This should only delete test_token_1
	deleteExpiredAccessTokens(suite.DB, suite.Client, suite.User)

	// Check the test_token_1 was deleted
	assert.True(
		suite.T(),
		suite.DB.Where(&AccessToken{Token: "test_token_1"}).First(&AccessToken{}).RecordNotFound(),
	)

	// Check the other three tokens are still around
	for _, token := range []string{"test_token_2", "test_token_3", "test_token_4"} {
		assert.False(
			suite.T(),
			suite.DB.Where(&AccessToken{Token: token}).First(&AccessToken{}).RecordNotFound(),
		)
	}

	// This should only delete test_token_2
	deleteExpiredAccessTokens(suite.DB, suite.Client, nil)

	// Check the test_token_2 was deleted
	assert.True(
		suite.T(),
		suite.DB.Where(&AccessToken{Token: "test_token_2"}).First(&AccessToken{}).RecordNotFound(),
	)

	// Check that last two tokens are still around
	for _, token := range []string{"test_token_3", "test_token_4"} {
		assert.False(
			suite.T(),
			suite.DB.Where(&AccessToken{Token: token}).First(&AccessToken{}).RecordNotFound(),
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
