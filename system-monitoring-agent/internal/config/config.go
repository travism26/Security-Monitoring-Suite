// internal/config/config.go
package config

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/spf13/viper"
)

// ConfigVersion tracks configuration changes
const CurrentConfigVersion = "1.0.0"

var (
	ErrInvalidTenantID     = errors.New("invalid tenant ID format")
	ErrMissingAPIKey       = errors.New("API key is required")
	ErrInvalidEndpoint     = errors.New("invalid endpoint URL")
	ErrInvalidRetention    = errors.New("invalid retention period")
	ErrInvalidStorageLimit = errors.New("invalid storage limit")
)

// TenantConfig holds tenant-specific configuration
type TenantConfig struct {
	ID        string `yaml:"ID"`
	APIKey    string `yaml:"APIKey"`
	Endpoints struct {
		Metrics       string `yaml:"Metrics"`
		HealthCheck   string `yaml:"HealthCheck"`
		KeyValidation string `yaml:"KeyValidation"`
	} `yaml:"Endpoints"`
}

// LogSettings holds logging configuration
type LogSettings struct {
	Level      string `yaml:"Level"`
	Format     string `yaml:"Format"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Brokers          []string `yaml:"Brokers"`
	Topic            string   `yaml:"Topic"`
	TenantTopic      string   `yaml:"TenantTopic"`
	SecurityProtocol string   `yaml:"SecurityProtocol"`
	SASLMechanism    string   `yaml:"SASLMechanism"`
	Username         string   `yaml:"Username"`
	Password         string   `yaml:"Password"`
}

// HTTPConfig holds HTTP-related configuration
type HTTPConfig struct {
	Endpoint      string `yaml:"Endpoint"`   // Maintaining backward compatibility
	StorageDir    string `yaml:"StorageDir"` // Maintaining backward compatibility
	RetryAttempts int    `yaml:"RetryAttempts"`
	RetryDelay    int    `yaml:"RetryDelay"`
	Timeout       int    `yaml:"Timeout"`
	Headers       struct {
		TenantID string `yaml:"TenantID"`
		APIKey   string `yaml:"APIKey"`
	} `yaml:"Headers"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	MaxStoragePerTenant int  `yaml:"MaxStoragePerTenant"`
	RetentionPeriod     int  `yaml:"RetentionPeriod"`
	CompressOldData     bool `yaml:"CompressOldData"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EncryptData bool     `yaml:"EncryptData"`
	ValidateSSL bool     `yaml:"ValidateSSL"`
	AllowedIPs  []string `yaml:"AllowedIPs"`
}

// Config represents the complete configuration structure
type Config struct {
	sync.RWMutex
	Version     string       `yaml:"Version"`
	Tenant      TenantConfig `yaml:"Tenant"`
	LogFilePath string       `yaml:"LogFilePath"`
	LogSettings LogSettings  `yaml:"LogSettings"`
	Interval    int          `yaml:"Interval"`
	Kafka       KafkaConfig  `yaml:"Kafka"`
	HTTP        HTTPConfig   `yaml:"HTTP"`
	StorageDir  string       `yaml:"StorageDir"` // Maintaining backward compatibility
	Monitors    struct {
		CPU     bool `yaml:"CPU"`
		Memory  bool `yaml:"Memory"`
		Disk    bool `yaml:"Disk"`
		Network bool `yaml:"Network"`
		Process bool `yaml:"Process"`
	} `yaml:"Monitors"`
	Thresholds struct {
		CPU                int `yaml:"CPU"`
		Memory             int `yaml:"Memory"`
		Disk               int `yaml:"Disk"`
		NetworkUtilization int `yaml:"NetworkUtilization"`
	} `yaml:"Thresholds"`
	Storage  StorageConfig  `yaml:"Storage"`
	Security SecurityConfig `yaml:"Security"`
}

// validateTenantID checks if the tenant ID matches the required format
func validateTenantID(id string) error {
	if id == "" {
		return ErrInvalidTenantID
	}
	// Tenant ID should be alphanumeric with optional hyphens, 4-32 chars
	match, _ := regexp.MatchString(`^[a-zA-Z0-9-]{4,32}$`, id)
	if !match {
		return ErrInvalidTenantID
	}
	return nil
}

// validateAPIKey checks if the API key is present and valid
func validateAPIKey(key string) error {
	if key == "" {
		return ErrMissingAPIKey
	}
	// API key should be at least 32 chars
	if len(key) < 32 {
		return ErrMissingAPIKey
	}
	return nil
}

// validateEndpoints checks if the endpoints are valid URLs
func validateEndpoints(cfg *Config) error {
	endpoints := []string{
		cfg.Tenant.Endpoints.Metrics,
		cfg.Tenant.Endpoints.HealthCheck,
		cfg.Tenant.Endpoints.KeyValidation,
	}

	urlPattern := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	for _, endpoint := range endpoints {
		if endpoint != "" && !urlPattern.MatchString(endpoint) {
			return fmt.Errorf("%w: %s", ErrInvalidEndpoint, endpoint)
		}
	}
	return nil
}

// validateStorage checks storage configuration
func validateStorage(cfg *Config) error {
	if cfg.Storage.RetentionPeriod < 1 {
		return ErrInvalidRetention
	}
	if cfg.Storage.MaxStoragePerTenant < 1 {
		return ErrInvalidStorageLimit
	}
	return nil
}

// Validate performs comprehensive configuration validation
func (cfg *Config) Validate() error {
	if err := validateTenantID(cfg.Tenant.ID); err != nil {
		return err
	}
	if err := validateAPIKey(cfg.Tenant.APIKey); err != nil {
		return err
	}
	if err := validateEndpoints(cfg); err != nil {
		return err
	}
	if err := validateStorage(cfg); err != nil {
		return err
	}
	return nil
}

// LoadConfig loads and validates the configuration
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Set default values
	setDefaultConfig()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config file found; using defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set current version
	cfg.Version = CurrentConfigVersion

	// Skip validation for empty/default config
	if cfg.Tenant.ID != "" || cfg.Tenant.APIKey != "" {
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid configuration: %w", err)
		}
	}

	return &cfg, nil
}

// setDefaultConfig sets default values for configuration
func setDefaultConfig() {
	viper.SetDefault("LogFilePath", "./agent.log")
	viper.SetDefault("LogSettings.Level", "info")
	viper.SetDefault("LogSettings.Format", "json")
	viper.SetDefault("LogSettings.MaxSize", 100)
	viper.SetDefault("LogSettings.MaxBackups", 3)
	viper.SetDefault("LogSettings.MaxAge", 28)
	viper.SetDefault("LogSettings.Compress", true)
	viper.SetDefault("Interval", 60)
	viper.SetDefault("Monitors.CPU", true)
	viper.SetDefault("Monitors.Memory", true)
	viper.SetDefault("Monitors.Disk", true)
	viper.SetDefault("Monitors.Network", true)
	viper.SetDefault("Monitors.Process", false)
	viper.SetDefault("Storage.MaxStoragePerTenant", 1024)
	viper.SetDefault("Storage.RetentionPeriod", 7)
	viper.SetDefault("Storage.CompressOldData", true)
	viper.SetDefault("Security.EncryptData", true)
	viper.SetDefault("Security.ValidateSSL", true)
	viper.SetDefault("Security.AllowedIPs", []string{})
}

// ReloadConfig reloads the configuration from disk
func (cfg *Config) ReloadConfig() error {
	cfg.Lock()
	defer cfg.Unlock()

	newCfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	// Update configuration
	*cfg = *newCfg

	return nil
}

// GetConfigVersion returns the current configuration version
func (cfg *Config) GetConfigVersion() string {
	cfg.RLock()
	defer cfg.RUnlock()
	return cfg.Version
}
