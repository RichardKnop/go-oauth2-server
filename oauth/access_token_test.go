package oauth

import (
	"fmt"
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

func TestGetClientIDUserIDQueryArgs(t *testing.T) {
	var queryParts, expectedQueryParts []string
	var args, expectedArgs []interface{}

	// client_id = 1 AND user_id IS NULL
	queryParts, args = getClientIDUserIDQueryArgs(&Client{ID: 1}, nil)
	expectedQueryParts = []string{"client_id = ?", "user_id IS NULL"}
	assert.Equal(
		t, expectedQueryParts, queryParts,
		"Query parts incorrect",
	)
	expectedArgs = []interface{}{uint(1)}
	assert.Equal(
		t, expectedArgs, args,
		"Args incorrect",
	)

	// client_id = 1 AND user_id = 2
	queryParts, args = getClientIDUserIDQueryArgs(&Client{ID: 1}, &User{ID: 2})
	expectedQueryParts = []string{"client_id = ?", "user_id = ?"}
	assert.Equal(
		t, expectedQueryParts, queryParts,
		"Query parts incorrect",
	)
	expectedArgs = []interface{}{uint(1), uint(2)}
	assert.Equal(
		t, expectedArgs, args,
		"Args incorrect",
	)
}

func (suite *OauthTestSuite) TestDeleteExpiredAccessTokens() {
	// Insert test access tokens that are expired
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_1",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_2",
		ExpiresAt: time.Now().Add(-10 * time.Second),
		Client:    *suite.Client,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}

	// Insert test access tokens that haven't expired yet
	if err := suite.DB.Create(&AccessToken{
		Token:     "test_token_3",
		ExpiresAt: time.Now().Add(+10 * time.Second),
		Client:    *suite.Client,
		User:      *suite.User,
		Scope:     "doesn't matter",
	}).Error; err != nil {
		log.Fatal(err)
	}
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
	assert.Equal(
		suite.T(),
		true,
		suite.DB.Where("token = ?", "test_token_1").First(&AccessToken{}).RecordNotFound(),
		"test_token_1 should be deleted",
	)

	// Check the other three tokens are still around
	for _, token := range []string{"test_token_2", "test_token_3", "test_token_4"} {
		assert.Equal(
			suite.T(),
			false,
			suite.DB.Where("token = ?", token).First(&AccessToken{}).RecordNotFound(),
			fmt.Sprintf("%s should still exist", token),
		)
	}

	// This should only delete test_token_2
	deleteExpiredAccessTokens(suite.DB, suite.Client, nil)

	// Check the test_token_2 was deleted
	assert.Equal(
		suite.T(),
		true,
		suite.DB.Where("token = ?", "test_token_2").First(&AccessToken{}).RecordNotFound(),
		"test_token_2 should be deleted",
	)

	// Check that last two tokens are still around
	for _, token := range []string{"test_token_3", "test_token_4"} {
		assert.Equal(
			suite.T(),
			false,
			suite.DB.Where("token = ?", token).First(&AccessToken{}).RecordNotFound(),
			fmt.Sprintf("%s should still exist", token),
		)
	}
}
