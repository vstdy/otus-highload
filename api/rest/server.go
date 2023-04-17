package rest

import (
	"net/http"

	"github.com/vstdy/otus-highload/service/project"
)

// NewServer returns new rest server.
func NewServer(svc project.Service, config Config) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router, err := NewRouter(svc, config)
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: config.ServerAddress, Handler: router}, nil
}
