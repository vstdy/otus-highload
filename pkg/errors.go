package pkg

import "errors"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrInvalidInput           = errors.New("invalid input")
	ErrNotFound               = errors.New("not found")
	ErrWrongCredentials       = errors.New("wrong credentials")
)
