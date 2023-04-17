//go:generate mockgen -source=interface.go -destination=./mock/service.go -package=servicemock
package project

import (
	"context"
	"github.com/google/uuid"
	"github.com/vstdy/otus-highload/model"
	"io"
)

type Service interface {
	io.Closer

	// CreateUser creates a new model.User.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)
	// GetUser returns user data.
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
}
