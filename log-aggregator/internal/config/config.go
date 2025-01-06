// internal/config/config.go
package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}
	API struct {
		Keys []string `mapstructure:"api_keys"`
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
	LogService struct {
		Environment string `mapstructure:"environment"`
		Application string `mapstructure:"application"`
		Component   string `mapstructure:"component"`
	}
}

// LoadConfig reads configuration from environment variables or config file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("api.api_keys", []string{"dev-api-key"}) // Default API key for development
	viper.SetDefault("kafka.topic", "logs")
	viper.SetDefault("kafka.groupid", "log-aggregator")
	viper.SetDefault("logservice.environment", "production")
	viper.SetDefault("logservice.application", "log-aggregator")
	viper.SetDefault("logservice.component", "log-service")

	// Map environment variables
	viper.SetEnvPrefix("LOG_AGG") // prefix for environment variables
	viper.AutomaticEnv()          // read environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Map specific environment variables
	viper.BindEnv("server.port", "LOG_AGG_SERVER_PORT")
	viper.BindEnv("server.host", "LOG_AGG_SERVER_HOST")
	viper.BindEnv("api.api_keys", "LOG_AGG_API_KEYS") // Comma-separated list of API keys
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.topic", "LOG_AGG_KAFKA_TOPIC")
	viper.BindEnv("kafka.group_id", "LOG_AGG_KAFKA_GROUP_ID")
	viper.BindEnv("database.host", "POSTGRES_HOST")
	viper.BindEnv("database.port", "POSTGRES_PORT")
	viper.BindEnv("database.user", "POSTGRES_USER")
	viper.BindEnv("database.password", "POSTGRES_PASSWORD")
	viper.BindEnv("database.name", "POSTGRES_DB")
	viper.BindEnv("logservice.environment", "LOG_AGG_ENV")
	viper.BindEnv("logservice.application", "LOG_AGG_APP")
	viper.BindEnv("logservice.component", "LOG_AGG_COMPONENT")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; using defaults and environment variables
		log.Println("No config file found. Using defaults and environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Process API keys from environment variable if present
	if apiKeys := viper.GetString("api.api_keys"); apiKeys != "" {
		config.API.Keys = strings.Split(apiKeys, ",")
		// Trim spaces from each key
		for i := range config.API.Keys {
			config.API.Keys[i] = strings.TrimSpace(config.API.Keys[i])
		}
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig ensures all required fields are set with valid values
func validateConfig(cfg *Config) error {
	if cfg.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if cfg.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}
	if len(cfg.Kafka.Brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker is required")
	}
	if cfg.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic is required")
	}
	if cfg.Kafka.GroupID == "" {
		return fmt.Errorf("kafka group ID is required")
	}
	if cfg.LogService.Environment == "" {
		return fmt.Errorf("log service environment is required")
	}
	if cfg.LogService.Application == "" {
		return fmt.Errorf("log service application name is required")
	}
	if cfg.LogService.Component == "" {
		return fmt.Errorf("log service component name is required")
	}
	if len(cfg.API.Keys) == 0 {
		return fmt.Errorf("at least one API key is required")
	}
	return nil
}
