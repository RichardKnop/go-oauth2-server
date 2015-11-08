package oauth

import (
	"errors"
	"strings"
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
	refreshToken, err := getOrCreateRefreshToken(cnf, db, client, user, scope)
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

func getClientIDUserIDQueryArgs(client *Client, user *User) ([]string, []interface{}) {
	// client_id
	queryParts := []string{"client_id = ?"}
	args := []interface{}{client.ID}

	// user_id
	if user == nil {
		queryParts = append(queryParts, "user_id IS NULL")
	} else {
		queryParts = append(queryParts, "user_id = ?")
		args = append(args, user.ID)
	}

	return queryParts, args
}

func deleteExpiredAccessTokens(db *gorm.DB, client *Client, user *User) {
	// Build a client/user id part of the query
	queryParts, args := getClientIDUserIDQueryArgs(client, user)

	// Add condition to query only for expired tokens
	queryParts = append(queryParts, "expires_at <= ?")
	args = append(args, time.Now())

	// And delete those tokens
	db.Where(strings.Join(queryParts, " AND "), args...).Delete(&AccessToken{})
}

func getOrCreateRefreshToken(cnf *config.Config, db *gorm.DB, client *Client, user *User, scope string) (*RefreshToken, error) {
	// Build a client/user id part of the query
	queryParts, args := getClientIDUserIDQueryArgs(client, user)

	// Try to fetch an existing refresh token first
	refreshToken := &RefreshToken{}
	notFound := db.Where(strings.Join(queryParts, " AND "), args...).First(refreshToken).RecordNotFound()

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
		refreshToken = newRefreshToken(cnf.RefreshTokenLifetime, client, user, scope)
		if err := db.Create(refreshToken).Error; err != nil {
			return nil, errors.New("Error saving refresh token")
		}
	}

	return refreshToken, nil
}

func respondWithAccessToken(w rest.ResponseWriter, cnf *config.Config, accessToken *AccessToken, refreshToken *RefreshToken) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.Token,
	})
}
