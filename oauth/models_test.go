package oauth

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	client := Client{AbstractModel: AbstractModel{ID: 1}}
	user := User{AbstractModel: AbstractModel{ID: 2}}

	var accessToken *AccessToken
	var value driver.Value
	var err error

	// When user object is nil
	accessToken = newAccessToken(
		3600,                   // expires in
		&client,                // client
		new(User),              // empty User object
		"scope doesn't matter", // scope
	)

	// accessToken.ClientID.Valid should be true
	assert.True(t, accessToken.ClientID.Valid)

	// accessToken.ClientID.Value() should return the object id, in this case int64(1)
	value, err = accessToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// accessToken.UserID.Valid should be false
	assert.False(t, accessToken.UserID.Valid)

	// accessToken.UserID.Value() should return nil
	value, err = accessToken.UserID.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When user object is not nil
	accessToken = newAccessToken(
		3600,    // expires in
		&client, // client
		&user,   // user
		"scope doesn't matter", // scope
	)

	// accessToken.ClientID.Valid should be true
	assert.True(t, accessToken.ClientID.Valid)

	// accessToken.ClientID.Value() should return the object id, in this case int64(1)
	value, err = accessToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// accessToken.UserID.Valid should be true
	assert.True(t, accessToken.UserID.Valid)

	// accessToken.UserID.Value() should return the object id, in this case int64(2)
	value, err = accessToken.UserID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), value)
}

func TestNewRefreshToken(t *testing.T) {
	client := Client{AbstractModel: AbstractModel{ID: 1}}
	user := User{AbstractModel: AbstractModel{ID: 2}}

	var refreshToken *RefreshToken
	var value driver.Value
	var err error

	// When user object is nil
	refreshToken = newRefreshToken(
		1209600,                // expires in
		&client,                // client
		new(User),              // empty User object
		"scope doesn't matter", // scope
	)

	// refreshToken.ClientID.Valid should be true
	assert.True(t, refreshToken.ClientID.Valid)

	// refreshToken.ClientID.Value() should return the object id, in this case int64(1)
	value, err = refreshToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// refreshToken.UserID.Valid should be false
	assert.False(t, refreshToken.UserID.Valid)

	// refreshToken.UserID.Value() should return nil
	value, err = refreshToken.UserID.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When user object is not nil
	refreshToken = newRefreshToken(
		1209600, // expires in
		&client, // client
		&user,   // user
		"scope doesn't matter", // scope
	)

	// accessToken.ClientID.Valid should be true
	assert.True(t, refreshToken.ClientID.Valid)

	// accessToken.ClientID.Value() should return the object id, in this case int64(1)
	value, err = refreshToken.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// refreshToken.UserID.Valid should be true
	assert.True(t, refreshToken.UserID.Valid)

	// refreshToken.UserID.Value() should return the object id, in this case int64(2)
	value, err = refreshToken.UserID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), value)
}

func TestNewAuthorizationCode(t *testing.T) {
	client := Client{AbstractModel: AbstractModel{ID: 1}}
	user := User{AbstractModel: AbstractModel{ID: 2}}

	var authorizationCode *AuthorizationCode
	var value driver.Value
	var err error

	// When user object is not nil
	authorizationCode = newAuthorizationCode(
		3600,    // expires in
		&client, // client
		&user,   // user
		"redirect URI doesn't matter", // redirect URI
		"scope doesn't matter",        // scope
	)

	// authorizationCode.ClientID.Valid should be true
	assert.True(t, authorizationCode.ClientID.Valid)

	// authorizationCode.ClientID.Value() should return the object id, in this case int64(1)
	value, err = authorizationCode.ClientID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)

	// authorizationCode.UserID.Valid should be true
	assert.True(t, authorizationCode.UserID.Valid)

	// authorizationCode.UserID.Value() should return the object id, in this case int64(2)
	value, err = authorizationCode.UserID.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(2), value)
}
