package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// TimestampModel ...
type TimestampModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// EmailTokenModel is an abstract model which can be used for objects from which
// we derive redirect emails (email confirmation, password reset and such)
type EmailTokenModel struct {
	gorm.Model
	Reference   string `sql:"type:varchar(40);unique;not null"`
	EmailSent   bool   `sql:"index;not null"`
	EmailSentAt *time.Time
	ExpiresAt   time.Time `sql:"index;not null"`
}
