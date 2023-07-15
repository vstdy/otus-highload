package common

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/api/rest"
	"github.com/vstdy/otus-highload/pkg"
	"github.com/vstdy/otus-highload/provider/broker/rabbitmq"
	"github.com/vstdy/otus-highload/provider/cache/redis"
	"github.com/vstdy/otus-highload/provider/config/etcd"
	"github.com/vstdy/otus-highload/service/project/v1"
	"github.com/vstdy/otus-highload/storage"
	"github.com/vstdy/otus-highload/storage/psql"
	"github.com/vstdy/otus-highload/storage/tarantool"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	Timeout          time.Duration    `mapstructure:"timeout"`
	LogLevel         zerolog.Level    `mapstructure:"-"`
	StorageType      string           `mapstructure:"storage_type"`
	HTTPServer       rest.Config      `mapstructure:"server,squash"`
	Cache            redis.Config     `mapstructure:"cache,squash"`
	Broker           rabbitmq.Config  `mapstructure:"broker,squash"`
	ExtConfig        etcd.Config      `mapstructure:"ext_config,squash"`
	PSQLStorage      psql.Config      `mapstructure:"psql_storage,squash"`
	TarantoolStorage tarantool.Config `mapstructure:"tarantool_storage,squash"`
}

const (
	psqlStorage = "psql"
)

// BuildDefaultConfig builds a Config with default values.
func BuildDefaultConfig() Config {
	return Config{
		Timeout:          5 * time.Second,
		LogLevel:         zerolog.InfoLevel,
		StorageType:      psqlStorage,
		HTTPServer:       rest.NewDefaultConfig(),
		Cache:            redis.NewDefaultConfig(),
		Broker:           rabbitmq.NewDefaultConfig(),
		ExtConfig:        etcd.NewDefaultConfig(),
		PSQLStorage:      psql.NewDefaultConfig(),
		TarantoolStorage: tarantool.NewDefaultConfig(),
	}
}

// BuildService builds project.Service dependency.
func (config Config) BuildService() (*project.Service, error) {
	cache, err := redis.NewClient(config.Cache)
	if err != nil {
		return nil, fmt.Errorf("building cache: %w", err)
	}

	broker, err := rabbitmq.NewClient(config.Broker)
	if err != nil {
		return nil, fmt.Errorf("building broker: %w", err)
	}

	extConfig, err := etcd.NewClient(config.ExtConfig)
	if err != nil {
		return nil, fmt.Errorf("building ext_config: %w", err)
	}

	st, err := config.BuildStorage()
	if err != nil {
		return nil, err
	}

	tr, err := tarantool.NewStorage(tarantool.WithConfig(config.TarantoolStorage))
	if err != nil {
		return nil, fmt.Errorf("building tarantool storage: %w", err)
	}

	svc, err := project.NewService(
		project.WithStorage(st),
		project.WithMessageStorage(tr),
		project.WithCache(cache),
		project.WithBroker(broker),
		project.WithExtConfig(extConfig),
	)
	if err != nil {
		return nil, fmt.Errorf("building service: %w", err)
	}

	return svc, nil
}

// BuildStorage builds storage dependency.
func (config Config) BuildStorage() (storage.IStorage, error) {
	var st storage.IStorage
	var err error

	switch config.StorageType {
	case psqlStorage:
		st, err = config.buildPsqlStorage()
	default:
		err = pkg.ErrUnsupportedStorageType
	}

	return st, err
}

// buildPsqlStorage builds psql.Storage dependency.
func (config Config) buildPsqlStorage() (*psql.Storage, error) {
	st, err := psql.NewStorage(
		psql.WithConfig(config.PSQLStorage),
	)
	if err != nil {
		return nil, fmt.Errorf("building psql storage: %w", err)
	}

	return st, nil
}
