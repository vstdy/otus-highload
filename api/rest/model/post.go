package model

import (
	"github.com/google/uuid"

	canonical "github.com/vstdy/otus-highload/model"
)

// CreatePostBody ...
type CreatePostBody struct {
	Text string `json:"text"`
}

// CreatePostResponse ...
type CreatePostResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

// NewCreatePostResponse ...
func NewCreatePostResponse(uuid uuid.UUID) CreatePostResponse {
	return CreatePostResponse{
		UUID: uuid,
	}
}

// UpdatePostBody ...
type UpdatePostBody struct {
	UUID uuid.UUID `json:"id"`
	Text string    `json:"text"`
}

func (up UpdatePostBody) ToCanonical() canonical.Post {
	return canonical.Post{
		UUID: up.UUID,
		Text: up.Text,
	}
}

// PostResponse ...
type PostResponse struct {
	UUID       uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	AuthorUUID uuid.UUID `json:"author_user_id"`
}

// NewPostResponse ...
func NewPostResponse(post canonical.PostExt) PostResponse {
	return PostResponse{
		UUID:       post.UUID,
		Text:       post.Text,
		AuthorUUID: post.AuthorUUID,
	}
}

// NewPostListResponse ...
func NewPostListResponse(posts []canonical.PostExt) []PostResponse {
	res := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		res = append(res, NewPostResponse(post))
	}

	return res
}
