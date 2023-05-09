package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/vstdy/otus-highload/pkg"
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
	DeletedAt  time.Time
}

// EncryptPassword ...
func (u User) EncryptPassword() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("encrypting password: %w", err)
	}

	return string(hash), nil
}

// ComparePasswords ...
func (u User) ComparePasswords(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return pkg.ErrWrongCredentials
	}

	return nil
}

// ValidateCreation ...
func (u User) ValidateCreation() error {
	if u.FirstName == "" {
		return fmt.Errorf("%w: first_name is empty", pkg.ErrInvalidInput)
	}
	if u.SecondName == "" {
		return fmt.Errorf("%w: second_name is empty", pkg.ErrInvalidInput)
	}
	if u.Age > 100 && u.Age < 7 {
		return fmt.Errorf("%w: age must be in range from 7 to 100", pkg.ErrInvalidInput)
	}
	if u.City == "" {
		return fmt.Errorf("%w: city is empty", pkg.ErrInvalidInput)
	}
	if u.Password == "" {
		return fmt.Errorf("%w: password is empty", pkg.ErrInvalidInput)
	}

	return nil
}

// ValidateAuthentication ...
func (u User) ValidateAuthentication() error {
	if u.ID < 1 {
		return fmt.Errorf("%w: id must be greater than 0", pkg.ErrInvalidInput)
	}
	if u.Password == "" {
		return fmt.Errorf("%w: password is empty", pkg.ErrInvalidInput)
	}

	return nil
}

// SearchUser ...
type SearchUser struct {
	FirstName string
	LastName  string
}
