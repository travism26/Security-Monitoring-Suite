// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string   `mapstructure:"topic"`
		GroupID string   `mapstructure:"group_id"`
	}
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
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
