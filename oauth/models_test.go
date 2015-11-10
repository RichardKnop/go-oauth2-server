package oauth

import (
	"testing"

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
	refreshToken = newRefreshToken(1209600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", refreshToken.Client.ClientID,
		"Refresh token should belong to test_client",
	)
	assert.Equal(
		t, uint(0), refreshToken.User.ID,
		"Refresh token should not belong to a user",
	)
}
