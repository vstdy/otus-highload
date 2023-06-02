package cache

import (
	"context"

	"github.com/go-redis/cache/v9"
)

type ICache interface {
	Set(item *cache.Item) error
	Once(item *cache.Item) error
	Get(ctx context.Context, key string, value interface{}) error
}
