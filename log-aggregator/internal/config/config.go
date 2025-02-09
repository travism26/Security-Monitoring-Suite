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
	Features struct {
		MultiTenancy struct {
			Enabled bool `mapstructure:"enabled"`
		} `mapstructure:"multi_tenancy"`
	} `mapstructure:"features"`
	Organization struct {
		ID   string `mapstructure:"id"`
		Name string `mapstructure:"name"`
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
		Host            string `mapstructure:"host"`
		Port            string `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		Name            string `mapstructure:"name"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"` // in minutes
		BatchSize       int    `mapstructure:"batch_size"`        // size of batch inserts
	}
	Cache struct {
		Enabled         bool `mapstructure:"enabled"`
		TTL             int  `mapstructure:"ttl"`              // in minutes
		TimeRangeTTL    int  `mapstructure:"time_range_ttl"`   // in minutes
		CleanupInterval int  `mapstructure:"cleanup_interval"` // in minutes
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
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.ttl", 5)              // 5 minutes default TTL
	viper.SetDefault("cache.time_range_ttl", 2)   // 2 minutes for time range queries
	viper.SetDefault("cache.cleanup_interval", 1) // 1 minute cleanup interval
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", 5)         // 5 minutes
	viper.SetDefault("database.batch_size", 1000)             // default batch size for inserts
	viper.SetDefault("api.api_keys", []string{"dev-api-key"}) // Default API key for development
	viper.SetDefault("kafka.topic", "logs")
	viper.SetDefault("kafka.groupid", "log-aggregator")
	viper.SetDefault("logservice.environment", "production")
	viper.SetDefault("logservice.application", "log-aggregator")
	viper.SetDefault("logservice.component", "log-service")
	viper.SetDefault("organization.id", "123e4567-e89b-12d3-a456-426614174000") // valid uuid
	viper.SetDefault("organization.name", "Default Organization")
	viper.SetDefault("features.multi_tenancy.enabled", false)

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
	viper.BindEnv("database.max_open_conns", "POSTGRES_MAX_OPEN_CONNS")
	viper.BindEnv("database.max_idle_conns", "POSTGRES_MAX_IDLE_CONNS")
	viper.BindEnv("database.conn_max_lifetime", "POSTGRES_CONN_MAX_LIFETIME")
	viper.BindEnv("database.batch_size", "POSTGRES_BATCH_SIZE")
	viper.BindEnv("cache.enabled", "LOG_AGG_CACHE_ENABLED")
	viper.BindEnv("cache.ttl", "LOG_AGG_CACHE_TTL")
	viper.BindEnv("cache.time_range_ttl", "LOG_AGG_CACHE_TIME_RANGE_TTL")
	viper.BindEnv("cache.cleanup_interval", "LOG_AGG_CACHE_CLEANUP_INTERVAL")
	viper.BindEnv("logservice.environment", "LOG_AGG_ENV")
	viper.BindEnv("logservice.application", "LOG_AGG_APP")
	viper.BindEnv("logservice.component", "LOG_AGG_COMPONENT")
	viper.BindEnv("organization.id", "LOG_AGG_ORG_ID")
	viper.BindEnv("organization.name", "LOG_AGG_ORG_NAME")
	viper.BindEnv("features.multi_tenancy.enabled", "LOG_AGG_MULTI_TENANCY_ENABLED")

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
	if cfg.Organization.ID == "" {
		return fmt.Errorf("organization ID is required")
	}
	if len(cfg.API.Keys) == 0 {
		return fmt.Errorf("at least one API key is required")
	}

	// Validate database connection pool settings
	if cfg.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("database max open connections must be greater than 0")
	}
	if cfg.Database.MaxIdleConns <= 0 {
		return fmt.Errorf("database max idle connections must be greater than 0")
	}
	if cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return fmt.Errorf("database max idle connections cannot be greater than max open connections")
	}
	if cfg.Database.ConnMaxLifetime <= 0 {
		return fmt.Errorf("database connection max lifetime must be greater than 0")
	}

	// Validate cache settings if enabled
	if cfg.Cache.Enabled {
		if cfg.Cache.TTL <= 0 {
			return fmt.Errorf("cache TTL must be greater than 0 when cache is enabled")
		}
		if cfg.Cache.TimeRangeTTL <= 0 {
			return fmt.Errorf("cache time range TTL must be greater than 0 when cache is enabled")
		}
		if cfg.Cache.CleanupInterval <= 0 {
			return fmt.Errorf("cache cleanup interval must be greater than 0 when cache is enabled")
		}
	}

	return nil
}
