package model

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// LoginBody ...
type LoginBody struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// ToCanonical ...
func (b LoginBody) ToCanonical() (model.User, error) {
	id, err := strconv.ParseInt(b.ID, 10, 64)
	if err != nil {
		return model.User{}, err
	}

	obj := model.User{
		ID:       id,
		Password: b.Password,
	}

	return obj, nil
}

// LoginResponse ...
type LoginResponse struct {
	Token uuid.UUID `json:"token"`
}

// NewLoginResponse ...
func NewLoginResponse(user model.User) LoginResponse {
	return LoginResponse{
		Token: user.UUID,
	}
}

// RegisterBody ...
type RegisterBody struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Age        uint8  `json:"age"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

// ToCanonical ...
func (b RegisterBody) ToCanonical() model.User {
	obj := model.User{
		FirstName:  b.FirstName,
		SecondName: b.SecondName,
		Age:        b.Age,
		Biography:  b.Biography,
		City:       b.City,
		Password:   b.Password,
	}

	return obj
}

// RegisterResponse ...
type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

// NewRegisterResponse ...
func NewRegisterResponse(user model.User) RegisterResponse {
	return RegisterResponse{
		UserID: user.UUID,
	}
}

// UserResponse ...
type UserResponse struct {
	UUID       uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Age        uint8     `json:"age"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
}

// NewGetUserResponse ...
func NewGetUserResponse(user model.User) UserResponse {
	return UserResponse{
		UUID:       user.UUID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Age:        user.Age,
		Biography:  user.Biography,
		City:       user.City,
	}
}

// NewSearchUsersResponse ...
func NewSearchUsersResponse(users []model.User) []UserResponse {
	res := make([]UserResponse, 0, len(users))
	for _, user := range users {
		obj := UserResponse{
			UUID:       user.UUID,
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			Age:        user.Age,
			Biography:  user.Biography,
			City:       user.City,
		}
		res = append(res, obj)
	}

	return res
}
