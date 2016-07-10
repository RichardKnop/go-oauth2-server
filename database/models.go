package database

import (
	"time"
)

type TimestampModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
