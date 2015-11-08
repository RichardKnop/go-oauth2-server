package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	client := Client{ClientID: "test_client"}
	user := User{Username: "test_username"}

	var accessToken *AccessToken

	accessToken = newAccessToken(3600, 3600, &client, &user, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, "test_username", accessToken.User.Username,
		"Access token should belong to test_username",
	)

	accessToken = newAccessToken(3600, 3600, &client, nil, "doesn't matter")
	assert.Equal(
		t, "test_client", accessToken.Client.ClientID,
		"Access token should belong to test_client",
	)
	assert.Equal(
		t, 0, accessToken.User.ID,
		"Access token should not belong to a user",
	)
}
