package handler

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"

	"github.com/vstdy/otus-highload/pkg/logging"
	"github.com/vstdy/otus-highload/service/project"
)

const (
	serviceName = "otus-project server"
)

// Handler keeps handler dependencies.
type Handler struct {
	service  project.IService
	jwtAuth  *jwtauth.JWTAuth
	logLevel zerolog.Level
}

// NewHandler returns a new Handler instance.
func NewHandler(service project.IService, jwtAuth *jwtauth.JWTAuth, logLevel zerolog.Level) Handler {
	return Handler{service: service, jwtAuth: jwtAuth, logLevel: logLevel}
}

// Logger returns logger with service field set.
func (h Handler) Logger(ctx context.Context) (context.Context, zerolog.Logger) {
	ctx, logger := logging.GetCtxLogger(ctx, logging.WithLogLevel(h.logLevel))
	logger = logger.With().Str(logging.ServiceKey, serviceName).Logger()

	return ctx, logger
}
