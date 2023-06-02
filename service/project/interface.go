//go:generate mockgen -source=interface.go -destination=./mock/service.go -package=servicemock
package project

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

type IService interface {
	io.Closer

	IUserService
	IFriendService
	IPostService
}

type IUserService interface {
	// CreateUser creates a new model.User.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)
	// GetUser returns user data.
	GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error)
	// SearchUsers searches users.
	SearchUsers(ctx context.Context, searchParams model.SearchUser) ([]model.User, error)
}

type IFriendService interface {
	// SetFriend adds friend to user.
	SetFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error
	// DeleteFriend deletes user's friend.
	DeleteFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error
}

type IPostService interface {
	// CreatePost creates post.
	CreatePost(ctx context.Context, userUUID uuid.UUID, text string) (uuid.UUID, error)
	// UpdatePost updates post.
	UpdatePost(ctx context.Context, userUUID uuid.UUID, post model.Post) error
	// DeletePost deletes post.
	DeletePost(ctx context.Context, userUUID, postUUID uuid.UUID) error
	// GetPost returns post.
	GetPost(ctx context.Context, postUUID uuid.UUID) (model.PostExt, error)
	// PostsFeed returns friends' most recent posts.
	PostsFeed(ctx context.Context, userUUID uuid.UUID, page model.Page) ([]model.PostExt, error)
}
