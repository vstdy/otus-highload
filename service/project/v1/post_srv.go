package project

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/google/uuid"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg"
)

const (
	defaultFeedCacheSize   = 1000
	defaultFeedCacheTTL    = time.Minute
	noCacheUsersConfigName = "no_cache_users"
)

// CreatePost creates post.
func (svc *Service) CreatePost(ctx context.Context, userUUID uuid.UUID, text string) (uuid.UUID, error) {
	if text == "" {
		return uuid.Nil, fmt.Errorf("%w: text is empty", pkg.ErrInvalidInput)
	}

	post, err := svc.storage.CreatePost(ctx, userUUID, text)
	if err != nil {
		return uuid.Nil, err
	}

	go svc.invalidateCache(post.AuthorID)

	go svc.publishNewPostNtf(post.ToPostExt(userUUID))

	return post.UUID, nil
}

// UpdatePost updates post.
func (svc *Service) UpdatePost(ctx context.Context, userUUID uuid.UUID, post model.Post) error {
	if post.Text == "" {
		return fmt.Errorf("%w: text is empty", pkg.ErrInvalidInput)
	}

	return svc.storage.UpdatePost(ctx, userUUID, post)
}

// DeletePost deletes post.
func (svc *Service) DeletePost(ctx context.Context, userUUID, postUUID uuid.UUID) error {
	return svc.storage.DeletePost(ctx, userUUID, postUUID)
}

// GetPost returns post.
func (svc *Service) GetPost(ctx context.Context, postUUID uuid.UUID) (model.PostExt, error) {
	return svc.storage.GetPost(ctx, postUUID)
}

// PostsFeed returns friends most recent posts.
func (svc *Service) PostsFeed(ctx context.Context, userUUID uuid.UUID, page model.Page) ([]model.PostExt, error) {
	if page.Offset+page.Limit < defaultFeedCacheSize {
		return svc.getFeedFromCache(ctx, userUUID, page)
	}

	return svc.storage.PostsFeed(ctx, userUUID, page)
}

// getFeedFromCache returns friends most recent posts from cache.
func (svc *Service) getFeedFromCache(ctx context.Context, userUUID uuid.UUID, page model.Page) ([]model.PostExt, error) {
	var posts []model.PostExt

	err := svc.cache.Once(&cache.Item{
		Key:   userUUID.String(),
		Value: &posts,
		TTL:   defaultFeedCacheTTL * 2,
		Do: func(*cache.Item) (interface{}, error) {
			return svc.storage.PostsFeed(ctx, userUUID, model.Page{Offset: 0, Limit: defaultFeedCacheSize})
		},
	})
	if err != nil {
		return nil, err
	}

	if page.Offset > len(posts) {
		return nil, nil
	}
	if page.Offset+page.Limit > len(posts) {
		return posts[page.Offset:], nil
	}
	return posts[page.Offset : page.Offset+page.Limit], nil
}

func (svc *Service) invalidateCache(userID int64) {
	response, err := svc.extConfig.Get(context.Background(), noCacheUsersConfigName)
	if err != nil {
		svc.Logger(nil).Warn().Err(err)
		return
	}

	if len(response.Kvs) == 0 {
		return
	}

	var noCacheUsers []int64
	err = json.Unmarshal(response.Kvs[0].Value, &noCacheUsers)
	if err != nil {
		svc.Logger(nil).Warn().Err(err).Msgf("unmarshal %s", noCacheUsersConfigName)
		return
	}

	for _, noCacheUser := range noCacheUsers {
		if userID != noCacheUser {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		setters, err := svc.storage.GetFriendSetters(ctx, userID)
		cancel()

		for _, u := range setters {
			err = svc.cache.Delete(context.Background(), u.String())
			if err != nil {
				svc.Logger(nil).Warn().Err(err)
			}
		}

		return
	}
}

func (svc *Service) publishNewPostNtf(post model.PostExt) {
	marshal, err := json.Marshal(post)
	if err != nil {
		svc.Logger(nil).Warn().Err(err)
	}

	err = svc.broker.Publish(marshal)
	if err != nil {
		svc.Logger(nil).Warn().Err(err)
	}
}
