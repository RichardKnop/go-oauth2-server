package migrations

import "time"

// Migration ...
type Migration struct {
	ID        int
	Name      string `sql:"size:255"`
	CreatedAt time.Time
}
