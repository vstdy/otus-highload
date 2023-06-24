package psql

import "fmt"

const (
	defaultURL = "postgres://user:password@localhost:5432/project?sslmode=disable"
)

// Config keeps Storage configuration.
type Config struct {
	URL             string `mapstructure:"database_url"`
	AsyncReplicaURL string `mapstructure:"async_replica_url"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.URL == "" {
		return fmt.Errorf("%s field: empty", "URL")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		URL:             defaultURL,
		AsyncReplicaURL: "",
	}
}
