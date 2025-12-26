package config

import (
	"errors"
	"fmt"
)

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	if err := c.Kafka.Validate(); err != nil {
		return fmt.Errorf("kafka config: %w", err)
	}

	if err := c.FCM.Validate(); err != nil {
		return fmt.Errorf("fcm config: %w", err)
	}

	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	return nil
}

// Validate validates server configuration
func (s *ServerConfig) Validate() error {
	if s.Port <= 0 || s.Port > 65535 {
		return errors.New("invalid server port")
	}
	return nil
}

// Validate validates database configuration
func (d *DatabaseConfig) Validate() error {
	if d.Host == "" {
		return errors.New("database host is required")
	}
	if d.Port <= 0 {
		return errors.New("database port is required")
	}
	if d.User == "" {
		return errors.New("database user is required")
	}
	if d.DBName == "" {
		return errors.New("database name is required")
	}
	return nil
}

// Validate validates Kafka configuration
func (k *KafkaConfig) Validate() error {
	if len(k.Brokers) == 0 {
		return errors.New("at least one Kafka broker is required")
	}
	if k.GroupID == "" {
		return errors.New("Kafka group ID is required")
	}
	if k.Topics.TaskCreated == "" {
		return errors.New("task created topic is required")
	}
	return nil
}

// Validate validates FCM configuration
func (f *FCMConfig) Validate() error {
	if f.CredentialsPath == "" {
		return errors.New("FCM credentials path is required")
	}
	if f.ProjectID == "" {
		return errors.New("FCM project ID is required")
	}
	return nil
}
