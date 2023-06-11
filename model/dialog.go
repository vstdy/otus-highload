package model

import (
	"time"

	"github.com/google/uuid"
)

// Dialog keeps dialog data.
type Dialog struct {
	UUID      uuid.UUID
	ChatID    uuid.UUID
	From      uuid.UUID
	To        uuid.UUID
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
