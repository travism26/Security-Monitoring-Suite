// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server struct {
		Port string
		Host string
	}
	Kafka struct {
		Brokers []string
		Topic   string
		GroupID string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
}

// LoadConfig reads configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("kafka.topic", "logs")
	viper.SetDefault("kafka.groupid", "log-aggregator")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
