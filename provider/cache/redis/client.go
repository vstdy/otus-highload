package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	cache *cache.Cache
}

// NewClient ...
func NewClient(config Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	options := &redis.Options{Addr: config.RedisAddress}

	rCache := cache.New(&cache.Options{
		Redis:      redis.NewClient(options),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &Client{cache: rCache}, nil
}

func (c *Client) Set(item *cache.Item) error {
	return c.cache.Set(item)
}

func (c *Client) Get(ctx context.Context, key string, value interface{}) error {
	return c.cache.Get(ctx, key, value)
}

func (c *Client) Once(item *cache.Item) error {
	return c.cache.Once(item)
}
