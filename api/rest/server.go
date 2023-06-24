package rest

import (
	"net/http"

	"github.com/vstdy/otus-highload/api/rest/hub"
	"github.com/vstdy/otus-highload/service/project"
)

// NewServer returns new rest server.
func NewServer(svc project.IService, hub *hub.Hub, config Config) (*http.Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	router, err := NewRouter(svc, hub, config)
	if err != nil {
		return nil, err
	}

	return &http.Server{Addr: config.ServerAddress, Handler: router}, nil
}
