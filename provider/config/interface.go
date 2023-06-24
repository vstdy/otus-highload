package ext_config

import (
	"context"
	"io"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type IExtConfig interface {
	io.Closer

	Put(ctx context.Context, key, val string) (*clientv3.PutResponse, error)
	Get(ctx context.Context, key string) (*clientv3.GetResponse, error)
	Delete(ctx context.Context, key string) (*clientv3.DeleteResponse, error)
}
