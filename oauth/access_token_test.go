package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	client := Client{ClientID: "test_client"}
	user := User{Username: "test_username"}

	var accessToken *AccessToken

	accessToken = newAccessToken(3600, &client, &user, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, "test_username", accessToken.User.Username,
		"Access token should belong to test_username",
	)

	accessToken = newAccessToken(3600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, 0, accessToken.User.ID,
		"Access token should not belong to a user",
	)
}

func TestNewRefreshToken(t *testing.T) {
	client := Client{ClientID: "test_client"}
	user := User{Username: "test_username"}

	var refreshToken *RefreshToken

	refreshToken = newRefreshToken(3600, &client, &user, "doesn't matter")
	assert.Equal(
		t, "test_client", refreshToken.Client.ClientID,
		"Refresh token should belong to test_client",
	)
	assert.Equal(
		t, "test_username", refreshToken.User.Username,
		"Refresh token should belong to test_username",
	)

	refreshToken = newRefreshToken(3600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", refreshToken.Client.ClientID,
		"Refresh token should belong to test_client",
	)
	assert.Equal(
		t, 0, refreshToken.User.ID,
		"Refresh token should not belong to a user",
	)
}

func TestGetClientIDUserIDQueryArgs(t *testing.T) {
	var queryParts, expectedQueryParts []string
	var args, expectedArgs []interface{}

	queryParts, args = getClientIDUserIDQueryArgs(&Client{ID: 1}, nil)
	expectedQueryParts = []string{"client_id = ?", "user_id IS NULL"}
	assert.Equal(
		t, expectedQueryParts, queryParts,
		"Query parts incorrect",
	)
	expectedArgs = []interface{}{1}
	assert.Equal(
		t, expectedArgs, args,
		"Args incorrect",
	)

	queryParts, args = getClientIDUserIDQueryArgs(&Client{ID: 1}, &User{ID: 2})
	expectedQueryParts = []string{"client_id = ?", "user_id = ?"}
	assert.Equal(
		t, expectedQueryParts, queryParts,
		"Query parts incorrect",
	)
	expectedArgs = []interface{}{1, 2}
	assert.Equal(
		t, expectedArgs, args,
		"Args incorrect",
	)
}

// func (suite *OauthTestSuite) TestDeleteExpiredAccessTokens() {
// 	// Insert test access tokens
// 	if err := suite.DB.Create(&AccessToken{
// 		Token:     "test_token_1",
// 		ExpiresAt: time.Now().Add(-10 * time.Second),
// 		Client:    *suite.Client,
// 		User:      *suite.User,
// 		Scope:     "doesn't matter",
// 	}).Error; err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := suite.DB.Create(&AccessToken{
// 		Token:     "test_token_2",
// 		ExpiresAt: time.Now().Add(-10 * time.Second),
// 		Client:    *suite.Client,
// 		Scope:     "doesn't matter",
// 	}).Error; err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := suite.DB.Create(&AccessToken{
// 		Token:     "test_token_3",
// 		ExpiresAt: time.Now().Add(+10 * time.Second),
// 		Client:    *suite.Client,
// 		User:      *suite.User,
// 		Scope:     "doesn't matter",
// 	}).Error; err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := suite.DB.Create(&AccessToken{
// 		Token:     "test_token_4",
// 		ExpiresAt: time.Now().Add(+10 * time.Second),
// 		Client:    *suite.Client,
// 		Scope:     "doesn't matter",
// 	}).Error; err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// This should only delete test_token_1
// 	deleteExpiredAccessTokens(suite.DB, suite.Client, suite.User)
// 	assert.Equal(
// 		suite.T(),
// 		true,
// 		suite.DB.Where(&AccessToken{Token: "test_token_1"}).First(&AccessToken{}).RecordNotFound(),
// 		"test_token_1 should be deleted",
// 	)
// 	for _, token := range []string{"test_token_2", "test_token_3", "test_token_4"} {
// 		assert.Equal(
// 			suite.T(),
// 			false,
// 			suite.DB.Where(&AccessToken{Token: token}).First(&AccessToken{}).RecordNotFound(),
// 			fmt.Sprintf("%s should still exist", token),
// 		)
// 	}
//
// 	// This should only delete test_token_2
// 	deleteExpiredAccessTokens(suite.DB, suite.Client, nil)
// 	assert.Equal(
// 		suite.T(),
// 		true,
// 		suite.DB.Where(&AccessToken{Token: "test_token_2"}).First(&AccessToken{}).RecordNotFound(),
// 		"test_token_2 should be deleted",
// 	)
// 	for _, token := range []string{"test_token_3", "test_token_4"} {
// 		assert.Equal(
// 			suite.T(),
// 			false,
// 			suite.DB.Where(&AccessToken{Token: token}).First(&AccessToken{}).RecordNotFound(),
// 			fmt.Sprintf("%s should still exist", token),
// 		)
// 	}
// }
