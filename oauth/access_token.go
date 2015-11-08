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

func grantAccessToken(cnf *config.Config, db *gorm.DB, client *Client, user *User, scope string) (*AccessToken, error) {
	// Fetch old access tokens for later deletion
	var oldAccessTokens []AccessToken
	queryParts := []string{"client_id = ?"}
	args := []interface{}{client.ID}
	if user != nil {
		queryParts = append(queryParts, "user_id = ?")
		args = append(args, user.ID)
	} else {
		queryParts = append(queryParts, "user_id IS NULL")
	}
	db.Where(strings.Join(queryParts, " AND "), args...).Preload("RefreshToken").Find(&oldAccessTokens)

	// Create a new access token
	accessToken := newAccessToken(cnf.AccessTokenLifetime, cnf.RefreshTokenLifetime, client, user, scope)
	if err := db.Create(accessToken).Error; err != nil {
		return nil, errors.New("Error saving access token")
	}

	// Delete old access tokens and their associated refresh tokens
	for _, oldAccessToken := range oldAccessTokens {
		db.Delete(&oldAccessToken)
		db.Delete(&oldAccessToken.RefreshToken)
	}

	return accessToken, nil
}

func newAccessToken(accessTokenLifetime int, refreshTokenLifetime int, client *Client, user *User, scope string) *AccessToken {
	accessToken := &AccessToken{
		AccessToken: uuid.New(),
		ExpiresAt:   time.Now().Add(time.Duration(accessTokenLifetime) * time.Second),
		Scope:       scope,
		Client:      *client,
		RefreshToken: RefreshToken{
			RefreshToken: uuid.New(),
			ExpiresAt:    time.Now().Add(time.Duration(refreshTokenLifetime) * time.Second),
		},
	}
	if user != nil {
		accessToken.User = *user
	}
	return accessToken
}

func respondWithAccessToken(w rest.ResponseWriter, cnf *config.Config, accessToken *AccessToken) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteJson(map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.AccessToken,
		"expires_in":    cnf.AccessTokenLifetime,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": accessToken.RefreshToken.RefreshToken,
	})
}
