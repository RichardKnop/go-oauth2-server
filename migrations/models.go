package migrations

import "time"

// Migration ...
type Migration struct {
	ID    int
	Name  string
	RunAt time.Time `sql:"not null"`
}
