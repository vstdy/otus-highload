package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ExtConfig struct {
	config Config
	etcd   *clientv3.Client
}

func NewClient(config Config) (*ExtConfig, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	client, err := clientv3.New(config.ToClientConfig())
	if err != nil {
		return nil, fmt.Errorf("etcd client: %w", err)
	}

	return &ExtConfig{config: config, etcd: client}, nil
}

func (b ExtConfig) Put(ctx context.Context, key, val string) (*clientv3.PutResponse, error) {
	return b.etcd.Put(ctx, key, val)
}

func (b ExtConfig) Get(ctx context.Context, key string) (*clientv3.GetResponse, error) {
	return b.etcd.Get(ctx, key)
}

func (b ExtConfig) Delete(ctx context.Context, key string) (*clientv3.DeleteResponse, error) {
	return b.etcd.Delete(ctx, key)
}

func (b ExtConfig) Close() error {
	return b.etcd.Close()
}
