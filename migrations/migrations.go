package migrations

import (
	"time"
)

// Migration represents a single database migration
type Migration struct {
	ID        int
	Name      string `sql:"size:255"`
	CreatedAt time.Time
}
