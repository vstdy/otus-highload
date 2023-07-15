package tarantool

import (
	"fmt"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

const (
	defaultAddress = "localhost:3301"
	defaultUser    = "user"
	defaultPass    = "password"
	defaultTimeout = 5 * time.Second
)

// Config keeps Storage configuration.
type Config struct {
	Address string        `mapstructure:"tarantool_address"`
	User    string        `mapstructure:"tarantool_user"`
	Pass    string        `mapstructure:"tarantool_pass"`
	Timeout time.Duration `mapstructure:"tarantool_timeout"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.Address == "" {
		return fmt.Errorf("%s field: empty", "tarantool_address")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Address: defaultAddress,
		User:    defaultUser,
		Pass:    defaultPass,
		//Timeout: defaultTimeout,
	}
}

func (config Config) ToTarantoolOpts() tarantool.Opts {
	return tarantool.Opts{
		User:    config.User,
		Pass:    config.Pass,
		Timeout: config.Timeout,
	}
}
