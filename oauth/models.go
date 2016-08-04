package oauth

import (
	"database/sql"
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
	"github.com/RichardKnop/uuid"
	"github.com/jinzhu/gorm"
)

// Client ...
type Client struct {
	gorm.Model
	Key         string         `sql:"type:varchar(254);unique;not null"`
	Secret      string         `sql:"type:varchar(60);not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (c *Client) TableName() string {
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
func (s *Scope) TableName() string {
	return "oauth_scopes"
}

// User ...
type User struct {
	gorm.Model
	Username   string         `sql:"type:varchar(254);unique;not null"`
	Password   sql.NullString `sql:"type:varchar(60)"`
	MetaUserID uint           `sql:"index"`
}

// TableName specifies table name
func (u *User) TableName() string {
	return "oauth_users"
}

// RefreshToken ...
type RefreshToken struct {
	gorm.Model
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
	Token     string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time `sql:"not null"`
	Scope     string    `sql:"type:varchar(200);not null"`
}

// TableName specifies table name
func (rt *RefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

// AccessToken ...
type AccessToken struct {
	gorm.Model
	ClientID  sql.NullInt64 `sql:"index;not null"`
	UserID    sql.NullInt64 `sql:"index"`
	Client    *Client
	User      *User
	Token     string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt time.Time `sql:"not null"`
	Scope     string    `sql:"type:varchar(200);not null"`
}

// TableName specifies table name
func (at *AccessToken) TableName() string {
	return "oauth_access_tokens"
}

// AuthorizationCode ...
type AuthorizationCode struct {
	gorm.Model
	ClientID    sql.NullInt64 `sql:"index;not null"`
	UserID      sql.NullInt64 `sql:"index;not null"`
	Client      *Client
	User        *User
	Code        string         `sql:"type:varchar(40);unique;not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
	ExpiresAt   time.Time      `sql:"not null"`
	Scope       string         `sql:"type:varchar(200);not null"`
}

// TableName specifies table name
func (ac *AuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

// NewRefreshToken creates new RefreshToken instance
func NewRefreshToken(client *Client, user *User, expiresIn int, scope string) *RefreshToken {
	refreshToken := &RefreshToken{
		ClientID:  util.PositiveIntOrNull(int64(client.ID)),
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		refreshToken.UserID = util.PositiveIntOrNull(int64(user.ID))
	}
	return refreshToken
}

// NewAccessToken creates new AccessToken instance
func NewAccessToken(client *Client, user *User, expiresIn int, scope string) *AccessToken {
	accessToken := &AccessToken{
		ClientID:  util.PositiveIntOrNull(int64(client.ID)),
		Token:     uuid.New(),
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		accessToken.UserID = util.PositiveIntOrNull(int64(user.ID))
	}
	return accessToken
}

// NewAuthorizationCode creates new AuthorizationCode instance
func NewAuthorizationCode(client *Client, user *User, expiresIn int, redirectURI, scope string) *AuthorizationCode {
	return &AuthorizationCode{
		ClientID:    util.PositiveIntOrNull(int64(client.ID)),
		UserID:      util.PositiveIntOrNull(int64(user.ID)),
		Code:        uuid.New(),
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
		RedirectURI: util.StringOrNull(redirectURI),
		Scope:       scope,
	}
}
