package rabbitmq

import (
	"fmt"
	"time"
)

const (
	defaultURL           = "amqp://user:password@localhost:5672"
	newPostsNtfQueueName = "new_posts_notification"
	publishDeadlineDur   = 5 * time.Second
)

// Config keeps Storage configuration.
type Config struct {
	RabbitmqURL                  string        `mapstructure:"rabbitmq_url"`
	RabbitmqNewPostsNtfQueueName string        `mapstructure:"rabbitmq_new_posts_ntf_queue_name"`
	RabbitmqPublishDeadlineDur   time.Duration `mapstructure:"rabbitmq_publish_deadline_dur"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.RabbitmqURL == "" {
		return fmt.Errorf("%s field: empty", "rabbitmq_url")
	}

	if config.RabbitmqNewPostsNtfQueueName == "" {
		return fmt.Errorf("%s field: empty", "rabbitmq_new_posts_ntf_queue_name")
	}

	if config.RabbitmqPublishDeadlineDur < time.Second {
		return fmt.Errorf("%s field: lower than 1 second", "rabbitmq_publish_deadline_dur")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		RabbitmqURL:                  defaultURL,
		RabbitmqNewPostsNtfQueueName: newPostsNtfQueueName,
		RabbitmqPublishDeadlineDur:   publishDeadlineDur,
	}
}
