package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
)

// SendDialog sends message to dialog.
func (svc *Service) SendDialog(ctx context.Context, from, to uuid.UUID, text string) error {
	if text == "" {
		return fmt.Errorf("%w: text is empty", pkg.ErrInvalidInput)
	}

	users, err := svc.storage.GetUsers(ctx, []uuid.UUID{from, to})
	if err != nil {
		return err
	}
	if len(users) != 2 {
		return pkg.ErrUserNotFound
	}

	chatID, err := svc.msgStorage.GetChat(ctx, users[0].ID, users[1].ID)
	if errors.Is(err, pkg.ErrNotFound) {
		chatID, err = svc.msgStorage.AddChat(ctx, users[0].ID, users[1].ID)
	}
	if err != nil {
		return err
	}

	return svc.msgStorage.SendDialog(ctx, chatID, from, to, text)
}

// ListDialog returns dialog messages.
func (svc *Service) ListDialog(ctx context.Context, from, to uuid.UUID, page model.Page) ([]model.Dialog, error) {
	users, err := svc.storage.GetUsers(ctx, []uuid.UUID{from, to})
	if err != nil {
		return nil, err
	}
	if len(users) != 2 {
		return nil, pkg.ErrUserNotFound
	}

	chatID, err := svc.msgStorage.GetChat(ctx, users[0].ID, users[1].ID)
	if errors.Is(err, pkg.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return svc.msgStorage.ListDialog(ctx, chatID, page)
}
