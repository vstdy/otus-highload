package migrations

import (
	"embed"
)

//go:embed *.sql
var migrations embed.FS

// GetMigrations returns migration files.
func GetMigrations() embed.FS {
	return migrations
}
