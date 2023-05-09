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
	Biography  sql.NullString
	City       string
	Password   sql.NullString
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
		Biography:  u.Biography.String,
		City:       u.City,
		Password:   u.Password.String,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeletedAt:  u.DeletedAt.Time,
	}
}

type Users []User

func (u Users) ToCanonical() []model.User {
	objs := make([]model.User, 0, len(u))
	for _, dbObj := range u {
		obj := dbObj.ToCanonical()
		objs = append(objs, obj)
	}

	return objs
}
