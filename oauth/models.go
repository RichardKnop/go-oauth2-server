package oauth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/RichardKnop/go-oauth2-server/json"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/pborman/uuid"
)

// Client ...
type Client struct {
	ID          int64          `gorm:"primary_key"`
	ClientID    string         `sql:"type:varchar(254);unique;not null"`
	Secret      string         `sql:"type:varchar(60);not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
}

// Scope ...
type Scope struct {
	ID          int64  `gorm:"primary_key"`
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description sql.NullString
	IsDefault   bool `sql:"default:false"`
}

// User ...
type User struct {
	ID       int64  `gorm:"primary_key"`
	Username string `sql:"type:varchar(254);unique;not null"`
	Password string `sql:"type:varchar(60);not null"`
}

// RefreshToken ...
type RefreshToken struct {
	ID        int64         `gorm:"primary_key"`
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
}

// AccessToken ...
type AccessToken struct {
	ID        int64         `gorm:"primary_key"`
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
}

// AuthorizationCode ...
type AuthorizationCode struct {
	ID          int64          `gorm:"primary_key"`
	Code        string         `sql:"type:varchar(40);unique;not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	ExpiresAt   time.Time      `sql:"not null"`
	Scope       string         `sql:"type:varchar(200);not null"`
	ClientID    sql.NullInt64  `sql:"index;not null"`
	UserID      sql.NullInt64  `sql:"index;not null"`
	Client      *Client
	User        *User
}

// newAccessToken creates new AccessToken instance
func newAccessToken(expiresIn int, client *Client, user *User, scope string) *AccessToken {
	clientID := util.IntOrNull(client.ID)
	userID := util.IntOrNull(user.ID)
	accessToken := &AccessToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
		ClientID:  clientID,
		UserID:    userID,
	}
	if clientID.Valid {
		accessToken.Client = client
	}
	if userID.Valid {
		accessToken.User = user
	}
	return accessToken
}

// newRefreshToken creates new RefreshToken instance
func newRefreshToken(expiresIn int, client *Client, user *User, scope string) *RefreshToken {
	clientID := util.IntOrNull(client.ID)
	userID := util.IntOrNull(user.ID)
	refreshToken := &RefreshToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
		ClientID:  clientID,
		UserID:    userID,
	}
	if clientID.Valid {
		refreshToken.Client = client
	}
	if userID.Valid {
		refreshToken.User = user
	}
	return refreshToken
}

// newAuthorizationCode creates new AuthorizationCode instance
func newAuthorizationCode(expiresIn int, client *Client, user *User, redirectURI, scope string) *AuthorizationCode {
	clientID := util.IntOrNull(client.ID)
	userID := util.IntOrNull(user.ID)
	authorizationCode := &AuthorizationCode{
		Code:        uuid.New(),
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
		RedirectURI: util.StringOrNull(redirectURI),
		Scope:       scope,
		ClientID:    clientID,
		UserID:      userID,
	}
	if clientID.Valid {
		authorizationCode.Client = client
	}
	if userID.Valid {
		authorizationCode.User = user
	}
	return authorizationCode
}

func writeJSON(w http.ResponseWriter, expiresIn int, accessToken *AccessToken, refreshToken *RefreshToken) {
	json.WriteJSON(w, map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    expiresIn,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.Token,
	}, 200)
}
