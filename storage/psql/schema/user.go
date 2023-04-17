package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// User keeps user data.
type User struct {
	ID         int64
	UUID       uuid.UUID
	FirstName  string
	SecondName string
	Age        uint8
	Biography  string
	City       string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  sql.NullTime
}

// ToCanonical converts a DB object to canonical model.
func (u User) ToCanonical() model.User {
	return model.User{
		ID:         u.ID,
		UUID:       u.UUID,
		FirstName:  u.FirstName,
		SecondName: u.SecondName,
		Age:        u.Age,
		Biography:  u.Biography,
		City:       u.City,
		Password:   u.Password,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeletedAt:  u.DeletedAt.Time,
	}
}
