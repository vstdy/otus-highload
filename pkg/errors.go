package pkg

import (
	"errors"
	"fmt"
)

const NotNullViolation = "23502"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrInvalidInput           = errors.New("invalid input")
	ErrNotFound               = errors.New("not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrWrongCredentials       = errors.New("wrong credentials")
)

type ErrInvalidUserArgs struct {
	Column string
}

func (err ErrInvalidUserArgs) Error() string {
	return fmt.Sprintf("'%s' column user not found", err.Column)
}
