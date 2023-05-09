package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
)

// CreateUser creates a new model.User.
func (svc *Service) CreateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	err := rawObj.ValidateCreation()
	if err != nil {
		return model.User{}, err
	}

	rawObj.Password, err = rawObj.EncryptPassword()
	if err != nil {
		return model.User{}, err
	}

	obj, err := svc.storage.CreateUser(ctx, rawObj)
	if err != nil {
		return model.User{}, fmt.Errorf("creating user: %w", err)
	}

	return obj, nil
}

// AuthenticateUser verifies the identity of credentials.
func (svc *Service) AuthenticateUser(ctx context.Context, rawObj model.User) (model.User, error) {
	err := rawObj.ValidateAuthentication()
	if err != nil {
		return model.User{}, err
	}

	obj, err := svc.storage.AuthenticateUser(ctx, rawObj)
	if err != nil {
		return model.User{}, fmt.Errorf("authenticating user: %w", err)
	}

	err = obj.ComparePasswords(rawObj.Password)
	if err != nil {
		return model.User{}, err
	}

	return obj, nil
}

// GetUser returns user data.
func (svc *Service) GetUser(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	return svc.storage.GetUser(ctx, userUUID)
}

// SearchUsers searches users.
func (svc *Service) SearchUsers(ctx context.Context, searchParams model.SearchUser) ([]model.User, error) {
	return svc.storage.SearchUsers(ctx, searchParams)
}
