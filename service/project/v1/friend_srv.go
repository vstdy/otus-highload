package project

import (
	"context"

	"github.com/google/uuid"
)

// SetFriend adds friend to user.
func (svc *Service) SetFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	return svc.storage.SetFriend(ctx, userUUID, friendUUID)
}

// DeleteFriend deletes user's friend.
func (svc *Service) DeleteFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	return svc.storage.DeleteFriend(ctx, userUUID, friendUUID)
}
