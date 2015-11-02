package service

import "time"

// Client ...
type Client struct {
	ID          int
	ClientID    string `sql:"type:varchar(254);unique;not null"`
	Password    string `sql:"type:varchar(60);not null"`
	RedirectURI string `sql:"type:varchar(200)"`
}

// Scope ...
type Scope struct {
	ID          int
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description string
	IsDefault   bool `sql:"default:false"`
}

// User ...
type User struct {
	ID       int
	Username string `sql:"type:varchar(254);unique;not null"`
	Password string `sql:"type:varchar(60);not null"`
}

// RefreshToken ...
type RefreshToken struct {
	ID           int
	RefreshToken string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt    time.Time `sql:"not null"`
}

// AccessToken ...
type AccessToken struct {
	ID             int
	AccessToken    string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt      time.Time `sql:"not null"`
	ClientID       int       `sql:"index;not null"`
	UserID         int       `sql:"index"`
	RefreshTokenID int       `sql:"index"`
	Scopes         []Scope   `gorm:"many2many:access_token_scopes;"`
	Client         Client
	User           User
	RefreshToken   RefreshToken
}

// AuthorizationCode ...
type AuthorizationCode struct {
	ID          int
	Code        string    `sql:"type:varchar(40);unique;not null"`
	RedirectURI string    `sql:"type:varchar(200)"`
	ExpiresAt   time.Time `sql:"not null"`
	ClientID    int       `sql:"index;not null"`
	UserID      int       `sql:"index"`
	Scopes      []Scope   `gorm:"many2many:access_token_scopes;"`
	Client      Client
	User        User
}
