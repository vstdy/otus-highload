package model

import (
	"time"

	"github.com/google/uuid"
)

// Post keeps post data.
type Post struct {
	UUID      uuid.UUID
	Text      string
	AuthorID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Post) ToPostExt(authorUUID uuid.UUID) PostExt {
	return PostExt{
		Post:       p,
		AuthorUUID: authorUUID,
	}
}

// PostExt keeps extended post data.
type PostExt struct {
	Post
	AuthorUUID uuid.UUID
}

func (p PostExt) ToNewPostNtf(users []uuid.UUID) NewPostNtf {
	return NewPostNtf{
		PostExt: p,
		Users:   users,
	}
}

// NewPostNtf keeps new post notification data.
type NewPostNtf struct {
	PostExt
	Users []uuid.UUID
}
