package database

import (
	"time"
)

// TimestampModel ...
type TimestampModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
