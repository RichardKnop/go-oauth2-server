package oauth

import (
	"errors"
	"time"

	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

func grantAccessToken(cnf *config.Config, db *gorm.DB, client *Client, user *User, scope string) (*AccessToken, *RefreshToken, error) {
	// Delete expired access tokens
	deleteExpiredAccessTokens(db, client, user)

	// Create a new access token
	accessToken := newAccessToken(cnf.AccessTokenLifetime, client, user, scope)
	if err := db.Create(accessToken).Error; err != nil {
		return nil, nil, errors.New("Error saving access token")
	}

	// Create or retrieve a refresh token
	refreshToken, err := getOrCreateRefreshToken(db, client, user, cnf.RefreshTokenLifetime, scope)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func newAccessToken(accessTokenLifetime int, client *Client, user *User, scope string) *AccessToken {
	accessToken := &AccessToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(accessTokenLifetime) * time.Second),
		Scope:     scope,
		Client:    *client,
	}
	if user != nil {
		accessToken.User = *user
	}
	return accessToken
}

func newRefreshToken(refreshTokenLifetime int, client *Client, user *User, scope string) *RefreshToken {
	refreshToken := &RefreshToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(refreshTokenLifetime) * time.Second),
		Scope:     scope,
		Client:    *client,
	}
	if user != nil {
		refreshToken.User = *user
	}
	return refreshToken
}

func deleteExpiredAccessTokens(db *gorm.DB, client *Client, user *User) {
	db.Where(&AccessToken{
		ClientID: clientIDOrNull(client),
		UserID:   userIDOrNull(user),
	}).Where("expires_at <= ?", time.Now()).Delete(&AccessToken{})
}

func getOrCreateRefreshToken(db *gorm.DB, client *Client, user *User, refreshTokenLifetime int, scope string) (*RefreshToken, error) {
	// Try to fetch an existing refresh token first
	refreshToken := &RefreshToken{}
	notFound := db.Where(&RefreshToken{
		ClientID: clientIDOrNull(client),
		UserID:   userIDOrNull(user),
	}).First(refreshToken).RecordNotFound()

	// Check if the token is expired, if found
	var expired bool
	if !notFound {
		expired = time.Now().After(refreshToken.ExpiresAt)
	}

	// If the refresh token has expired, delete it
	if expired {
		db.Delete(refreshToken)
	}

	// Create a new refresh token if it expired or was not found
	if expired || notFound {
		refreshToken = newRefreshToken(refreshTokenLifetime, client, user, scope)
		if err := db.Create(refreshToken).Error; err != nil {
			return nil, errors.New("Error saving refresh token")
		}
	}

	return refreshToken, nil
}

func respondWithAccessToken(w rest.ResponseWriter, cnf *config.Config, accessToken *AccessToken, refreshToken *RefreshToken) {
	// Content-Type header must set charset in response
	// See https://github.com/ant0ine/go-json-rest/issues/156
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write access token to response
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.Token,
	})
}
