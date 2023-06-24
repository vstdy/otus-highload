package project

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg/logging"
)

// SetFriend adds friend to user.
func (svc *Service) SetFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	return svc.storage.SetFriend(ctx, userUUID, friendUUID)
}

// DeleteFriend deletes user's friend.
func (svc *Service) DeleteFriend(ctx context.Context, userUUID, friendUUID uuid.UUID) error {
	return svc.storage.DeleteFriend(ctx, userUUID, friendUUID)
}

// consumeNewPostsNotifications ...
func (svc *Service) consumeNewPostsNotifications() {
	logger := svc.Logger(nil).With().Str(logging.JobKey, newPostsNtfJobName).Logger()

	msgs, err := svc.broker.Consume()
	if err != nil {
		logger.Err(err).Msg("failed to register a consumer")
		return
	}

	for msg := range msgs {
		var post model.PostExt
		err = json.Unmarshal(msg.Body, &post)
		if err != nil {
			logger.Err(err).Msg("failed to unmarshal post")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		setters, err := svc.storage.GetFriendSetters(ctx, post.AuthorID)
		cancel()
		if err != nil {
			logger.Err(err).Msg("failed to fetch setters")
			continue
		}
		notification := model.NewPostNtf{
			PostExt: post,
			Users:   setters,
		}
		svc.hub <- notification
	}

	logger.Warn().Msg("exiting consumer")
}
