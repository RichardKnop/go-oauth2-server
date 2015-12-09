package oauth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/RichardKnop/go-oauth2-server/response"
	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

// Client ...
type Client struct {
	gorm.Model
	ClientID    string         `sql:"type:varchar(254);unique;not null"`
	Secret      string         `sql:"type:varchar(60);not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (c Client) TableName() string {
	return "oauth_clients"
}

// Scope ...
type Scope struct {
	gorm.Model
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description sql.NullString
	IsDefault   bool `sql:"default:false"`
}

// TableName specifies table name
func (s Scope) TableName() string {
	return "oauth_scopes"
}

// User ...
type User struct {
	gorm.Model
	Username string `sql:"type:varchar(254);unique;not null"`
	Password string `sql:"type:varchar(60);not null"`
}

// TableName specifies table name
func (u User) TableName() string {
	return "oauth_users"
}

// RefreshToken ...
type RefreshToken struct {
	gorm.Model
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
}

// TableName specifies table name
func (rt RefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

// AccessToken ...
type AccessToken struct {
	gorm.Model
	Token     string        `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time     `sql:"not null"`
	Scope     string        `sql:"type:varchar(200);not null"`
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
}

// TableName specifies table name
func (at AccessToken) TableName() string {
	return "oauth_access_tokens"
}

// AuthorizationCode ...
type AuthorizationCode struct {
	gorm.Model
	Code        string         `sql:"type:varchar(40);unique;not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	ExpiresAt   time.Time      `sql:"not null"`
	Scope       string         `sql:"type:varchar(200);not null"`
	ClientID    sql.NullInt64  `sql:"index;not null"`
	UserID      sql.NullInt64  `sql:"index;not null"`
	Client      *Client
	User        *User
}

// TableName specifies table name
func (ac AuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
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
	response.WriteJSON(w, map[string]interface{}{
		"id":            accessToken.ID,
		"access_token":  accessToken.Token,
		"expires_in":    expiresIn,
		"token_type":    "Bearer",
		"scope":         accessToken.Scope,
		"refresh_token": refreshToken.Token,
	}, 200)
}
