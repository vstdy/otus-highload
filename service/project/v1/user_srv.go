package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
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
	user, err := svc.storage.GetUsers(ctx, []uuid.UUID{userUUID})
	if err != nil {
		return model.User{}, err
	}
	if len(user) != 1 {
		return model.User{}, pkg.ErrNotFound
	}

	return user[0], nil
}

// SearchUsers searches users.
func (svc *Service) SearchUsers(ctx context.Context, firstName, secondName string) ([]model.User, error) {
	return svc.storage.SearchUsers(ctx, firstName, secondName)
}
