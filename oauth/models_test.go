package oauth

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	client := Client{ID: 1}
	user := User{ID: 2}

	var accessToken *AccessToken
	var value driver.Value
	var err error

	// With user
	accessToken = newAccessToken(3600, &client, &user, "doesn't matter")

	assert.True(t, accessToken.ClientID.Valid)
	value, err = accessToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	assert.True(t, accessToken.UserID.Valid)
	value, err = accessToken.UserID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), value)

	// Without user
	accessToken = newAccessToken(3600, &client, nil, "doesn't matter")

	assert.True(t, accessToken.ClientID.Valid)
	value, err = accessToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	assert.False(t, accessToken.UserID.Valid)
}

func TestNewRefreshToken(t *testing.T) {
	client := Client{ID: 1}
	user := User{ID: 2}

	var refreshToken *RefreshToken
	var value driver.Value
	var err error

	// With user
	refreshToken = newRefreshToken(3600, &client, &user, "doesn't matter")

	assert.True(t, refreshToken.ClientID.Valid)
	value, err = refreshToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	assert.True(t, refreshToken.UserID.Valid)
	value, err = refreshToken.UserID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), value)

	// Without user
	refreshToken = newRefreshToken(1209600, &client, nil, "doesn't matter")

	assert.True(t, refreshToken.ClientID.Valid)
	value, err = refreshToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	assert.False(t, refreshToken.UserID.Valid)
}
