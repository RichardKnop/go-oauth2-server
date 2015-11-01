package service

import "time"

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
	UserID         int       `sql:"index"`
	User           User
	RefreshTokenID int `sql:"index"`
	RefreshToken   RefreshToken
}
