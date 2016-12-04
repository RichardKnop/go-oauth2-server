package migrations

import (
	"github.com/jinzhu/gorm"
)

// Migration represents a single database migration
type Migration struct {
	gorm.Model
	Name string `sql:"size:255"`
}
