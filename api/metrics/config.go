package metrics

import (
	"fmt"
)

// Config keeps rest params.
type Config struct {
	MetricsServerAddress string `mapstructure:"metrics_server_address"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {

	if config.MetricsServerAddress == "" {
		return fmt.Errorf("%s field: empty", "metrics_server_address")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		MetricsServerAddress: "0.0.0.0:9100",
	}
}
