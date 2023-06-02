package common

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/api/metrics"
	"github.com/vstdy/otus-highload/api/rest"
	"github.com/vstdy/otus-highload/pkg"
	"github.com/vstdy/otus-highload/provider/cache/redis"
	"github.com/vstdy/otus-highload/service/project/v1"
	"github.com/vstdy/otus-highload/storage"
	"github.com/vstdy/otus-highload/storage/psql"
)

// Config combines sub-configs for all services, storages and providers.
type Config struct {
	Timeout           time.Duration  `mapstructure:"timeout"`
	LogLevel          zerolog.Level  `mapstructure:"-"`
	StorageType       string         `mapstructure:"storage_type"`
	HTTPServer        rest.Config    `mapstructure:"server,squash"`
	HTTPMetricsServer metrics.Config `mapstructure:"metrics_server,squash"`
	Cache             redis.Config   `mapstructure:"cache,squash"`
	PSQLStorage       psql.Config    `mapstructure:"psql_storage,squash"`
}

const (
	psqlStorage = "psql"
)

// BuildDefaultConfig builds a Config with default values.
func BuildDefaultConfig() Config {
	return Config{
		Timeout:           5 * time.Second,
		LogLevel:          zerolog.InfoLevel,
		StorageType:       psqlStorage,
		HTTPServer:        rest.NewDefaultConfig(),
		HTTPMetricsServer: metrics.NewDefaultConfig(),
		Cache:             redis.NewDefaultConfig(),
		PSQLStorage:       psql.NewDefaultConfig(),
	}
}

// BuildService builds project.Service dependency.
func (config Config) BuildService() (*project.Service, error) {
	cache, err := redis.NewClient(config.Cache)
	if err != nil {
		return nil, fmt.Errorf("building cache: %w", err)
	}

	st, err := config.BuildStorage()
	if err != nil {
		return nil, err
	}

	svc, err := project.NewService(
		project.WithStorage(st),
		project.WithCache(cache),
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
