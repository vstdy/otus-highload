package project

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/pkg/logging"
	"github.com/vstdy/otus-highload/service/project"
	inter "github.com/vstdy/otus-highload/storage"
)

const (
	serviceName = "Project service"
)

var _ project.Service = (*Service)(nil)

type (
	// Service keeps service dependencies.
	Service struct {
		storage inter.Storage
	}

	// ServiceOption defines functional argument for Service constructor.
	ServiceOption func(*Service) error
)

// WithStorage sets Storage.
func WithStorage(st inter.Storage) ServiceOption {
	return func(svc *Service) error {
		svc.storage = st

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

	return svc, nil
}

// Close closes all service dependencies.
func (svc *Service) Close() error {
	if svc.storage == nil {
		return nil
	}

	if err := svc.storage.Close(); err != nil {
		return fmt.Errorf("closing storage: %w", err)
	}

	return nil
}

// Logger returns logger with service field set.
func (svc *Service) Logger(ctx context.Context) *zerolog.Logger {
	_, logger := logging.GetCtxLogger(ctx)
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return &logger
}
