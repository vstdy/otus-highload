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
	ErrWrongCredentials       = errors.New("wrong credentials")
)

type ErrSetFriendInvalidArgs struct {
	Column string
}

func (err ErrSetFriendInvalidArgs) Error() string {
	return fmt.Sprintf("%s not found", err.Column)
}
