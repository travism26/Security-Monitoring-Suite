package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Environment represents the runtime environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// Config represents the application configuration
type Config struct {
	Environment Environment `yaml:"environment"`
	Endpoint    string      `yaml:"endpoint"`
	Interval    int         `yaml:"interval"`
	LogLevel    string      `yaml:"log_level"`
	Scenarios   struct {
		HighCPU struct {
			CPUUsage float64 `yaml:"cpu_usage"`
			Duration int     `yaml:"duration"`
		} `yaml:"high-cpu"`
		MaliciousProcess struct {
			ProcessName   string  `yaml:"process_name"`
			CPUPercent    float64 `yaml:"cpu_percent"`
			MemoryPercent float64 `yaml:"memory_percent"`
		} `yaml:"malicious-process"`
	} `yaml:"scenarios"`
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate environment
	switch c.Environment {
	case Development, Staging, Production:
		// Valid environment
	case "":
		c.Environment = Development // Default to development
	default:
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	// Validate endpoint
	if c.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}

	// Validate interval
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be greater than 0")
	}

	// Validate log level
	switch strings.ToUpper(c.LogLevel) {
	case "DEBUG", "INFO", "WARN", "ERROR":
		// Valid log level
	case "":
		c.LogLevel = "INFO" // Default to INFO
	default:
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Validate scenarios
	if c.Scenarios.HighCPU.CPUUsage <= 0 || c.Scenarios.HighCPU.CPUUsage > 100 {
		return fmt.Errorf("high CPU usage must be between 0 and 100")
	}
	if c.Scenarios.HighCPU.Duration <= 0 {
		return fmt.Errorf("high CPU duration must be greater than 0")
	}

	if c.Scenarios.MaliciousProcess.ProcessName == "" {
		return fmt.Errorf("malicious process name is required")
	}
	if c.Scenarios.MaliciousProcess.CPUPercent < 0 || c.Scenarios.MaliciousProcess.CPUPercent > 100 {
		return fmt.Errorf("malicious process CPU percent must be between 0 and 100")
	}
	if c.Scenarios.MaliciousProcess.MemoryPercent < 0 || c.Scenarios.MaliciousProcess.MemoryPercent > 100 {
		return fmt.Errorf("malicious process memory percent must be between 0 and 100")
	}

	return nil
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	// Get environment from ENV var or default to development
	env := Environment(strings.ToLower(os.Getenv("APP_ENV")))
	if env == "" {
		env = Development
	}

	// Construct environment-specific config path
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	envPath := filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, env, ext))

	// Try to load environment-specific config first, fall back to default
	data, err := os.ReadFile(envPath)
	if err != nil {
		// If env-specific config doesn't exist, try default config
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set environment from ENV var
	config.Environment = env

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}
