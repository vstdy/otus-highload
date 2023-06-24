package etcd

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultDialTimeout = 5 * time.Second

var (
	defaultEtcdEndpoits = []string{"localhost:2379"}
)

// Config keeps ExtConfig configuration.
type Config struct {
	EtcdEndpoints []string      `mapstructure:"etcd_endpoints"`
	DialTimeout   time.Duration `mapstructure:"dial_timeout"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if len(config.EtcdEndpoints) == 0 {
		return fmt.Errorf("%s field: empty", "etcd_endpoints")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		EtcdEndpoints: defaultEtcdEndpoits,
		DialTimeout:   defaultDialTimeout,
	}
}

func (config Config) ToClientConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints:   config.EtcdEndpoints,
		DialTimeout: config.DialTimeout,
	}
}
