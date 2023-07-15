package migrations

import (
	"embed"
)

//go:embed *.lua
var migrations embed.FS

// GetMigrations returns migration files.
func GetMigrations() embed.FS {
	return migrations
}
