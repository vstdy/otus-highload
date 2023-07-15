//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// IStorage ...
type IStorage interface {
	io.Closer

	IMigration
	IGenerate

	IUserStorage
	IFriendStorage
	IPostStorage
	IMessageStorage
}

type IMigration interface {
	// MigrateUp applies all available migrations.
	MigrateUp() error
	// MigrateDown rolls back a single migration from the current version.
	MigrateDown() error
}

type IGenerate interface {
	// CopyUsers copies users to storage.
	CopyUsers(ctx context.Context, objs []model.User) (int64, error)
	// CopyFriends copies friends to storage.
	CopyFriends(ctx context.Context, objs []model.Friend) (int64, error)
	// CopyPosts copies posts to storage.
	CopyPosts(ctx context.Context, objs []model.Post) (int64, error)
}

type IUserStorage interface {
	// CreateUser adds given objects to storage.
	CreateUser(ctx context.Context, obj model.User) (model.User, error)
	// AuthenticateUser verifies the identity of credentials.
	AuthenticateUser(ctx context.Context, obj model.User) (model.User, error)
	// GetUsers returns users data.
	GetUsers(ctx context.Context, userUUIDs []uuid.UUID) ([]model.User, error)
	// SearchUsers searches users.
	SearchUsers(ctx context.Context, firstName, secondName string) ([]model.User, error)
}

type IFriendStorage interface {
	// SetFriend adds friend to user.
	SetFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error
	// DeleteFriend deletes user's friend.
	DeleteFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error
	// GetFriendSetters returns users who added given user to friends.
	GetFriendSetters(ctx context.Context, friendID int64) ([]uuid.UUID, error)
}

type IPostStorage interface {
	// CreatePost creates post.
	CreatePost(ctx context.Context, userUUID uuid.UUID, text string) (model.Post, error)
	// UpdatePost updates post.
	UpdatePost(ctx context.Context, userUUID uuid.UUID, post model.Post) error
	// DeletePost deletes post.
	DeletePost(ctx context.Context, userUUID, postUUID uuid.UUID) error
	// GetPost returns post.
	GetPost(ctx context.Context, postUUID uuid.UUID) (model.PostExt, error)
	// PostsFeed returns friends' most recent posts.
	PostsFeed(ctx context.Context, userUUID uuid.UUID, page model.Page) ([]model.PostExt, error)
}

type IMessageStorage interface {
	IChatStorage
	IDialogStorage
}

type IChatStorage interface {
	// AddChat adds new chat.
	AddChat(ctx context.Context, user1, user2 int64) (uuid.UUID, error)
	// GetChat returns chat.
	GetChat(ctx context.Context, user1, user2 int64) (uuid.UUID, error)
}

type IDialogStorage interface {
	// SendDialog sends message to dialog.
	SendDialog(ctx context.Context, chatID, from, to uuid.UUID, text string) error
	// ListDialog returns dialog messages.
	ListDialog(ctx context.Context, chatID uuid.UUID, page model.Page) ([]model.Dialog, error)
}
