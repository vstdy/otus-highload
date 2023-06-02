package redis

import "fmt"

const (
	redisAddressDefault = "localhost:6379"
)

// Config keeps Storage configuration.
type Config struct {
	RedisAddress string `mapstructure:"redis_address"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.RedisAddress == "" {
		return fmt.Errorf("%s field: empty", "redis_address")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		RedisAddress: redisAddressDefault,
	}
}
