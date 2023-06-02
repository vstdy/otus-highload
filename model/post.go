package model

import (
	"time"

	"github.com/google/uuid"
)

// Post keeps post data.
type Post struct {
	ID        int64
	UUID      uuid.UUID
	Text      string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostExt keeps extended post data.
type PostExt struct {
	Post
	AuthorUUID uuid.UUID
}
