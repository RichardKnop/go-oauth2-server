package migrations

import "time"

// Migration ...
type Migration struct {
	ID        int
	Name      string
	CreatedAt time.Time `sql:"not null"`
}
