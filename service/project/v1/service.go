package project

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/model"
	"github.com/vstdy/otus-highload/pkg/logging"
	"github.com/vstdy/otus-highload/provider/broker"
	"github.com/vstdy/otus-highload/provider/cache"
	ext_config "github.com/vstdy/otus-highload/provider/config"
	"github.com/vstdy/otus-highload/service/project"
	"github.com/vstdy/otus-highload/storage"
)

const (
	serviceName        = "Project service"
	newPostsNtfJobName = "New posts notifications"
)

var _ project.IService = (*Service)(nil)

type (
	// Service keeps service dependencies.
	Service struct {
		storage    storage.IStorage
		msgStorage storage.IMessageStorage
		cache      cache.ICache
		broker     broker.IBroker
		extConfig  ext_config.IExtConfig
		hub        chan model.NewPostNtf
	}

	// ServiceOption defines functional argument for Service constructor.
	ServiceOption func(*Service) error
)

// WithStorage sets IStorage.
func WithStorage(st storage.IStorage) ServiceOption {
	return func(svc *Service) error {
		svc.storage = st

		return nil
	}
}

// WithMessageStorage sets IMessageStorage.
func WithMessageStorage(st storage.IMessageStorage) ServiceOption {
	return func(svc *Service) error {
		svc.msgStorage = st

		return nil
	}
}

// WithCache sets ICache.
func WithCache(c cache.ICache) ServiceOption {
	return func(svc *Service) error {
		svc.cache = c

		return nil
	}
}

// WithBroker sets IBroker.
func WithBroker(b broker.IBroker) ServiceOption {
	return func(svc *Service) error {
		svc.broker = b

		return nil
	}
}

// WithExtConfig sets IExtConfig.
func WithExtConfig(c ext_config.IExtConfig) ServiceOption {
	return func(svc *Service) error {
		svc.extConfig = c

		return nil
	}
}

// NewService creates a new project service.
func NewService(opts ...ServiceOption) (*Service, error) {
	svc := new(Service)
	for optIdx, opt := range opts {
		if err := opt(svc); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if svc.storage == nil {
		return nil, fmt.Errorf("storage: nil")
	}

	svc.hub = make(chan model.NewPostNtf)
	go svc.consumeNewPostsNotifications()

	return svc, nil
}

// Close closes all service dependencies.
func (svc *Service) Close() error {
	if err := svc.broker.Close(); err != nil {
		return fmt.Errorf("closing broker: %w", err)
	}

	if err := svc.storage.Close(); err != nil {
		return fmt.Errorf("closing storage: %w", err)
	}

	close(svc.hub)

	return nil
}

// Logger returns logger with service field set.
func (svc *Service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}

func (svc *Service) GetHub() chan model.NewPostNtf {
	return svc.hub
}
