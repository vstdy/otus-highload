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

	chatID, err := svc.storage.GetChat(ctx, from, to)
	if errors.Is(err, pkg.ErrNotFound) {
		chatID, err = svc.storage.AddChat(ctx, from, to)
	}
	if err != nil {
		return err
	}

	return svc.storage.SendDialog(ctx, chatID, from, to, text)
}

// ListDialog returns dialog messages.
func (svc *Service) ListDialog(ctx context.Context, from, to uuid.UUID, page model.Page) ([]model.Dialog, error) {
	chatID, err := svc.storage.GetChat(ctx, from, to)
	if errors.Is(err, pkg.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return svc.storage.ListDialog(ctx, chatID, page)
}
