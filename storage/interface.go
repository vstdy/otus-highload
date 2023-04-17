//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// Storage ...
type Storage interface {
	io.Closer
	Migrate
	User
}

type Migrate interface {
	// MigrateUp applies all available migrations.
	MigrateUp() error
	// MigrateDown rolls back a single migration from the current version.
	MigrateDown() error
}

type User interface {
	// CreateUser adds given objects to storage.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)
	// GetUser returns user data.
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
