package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		configData  string
		expectError bool
		expected    Config
	}{
		{
			name:       "Valid tenant config file",
			configFile: "config.yaml",
			configData: `
Version: "1.0.0"
Tenant:
  ID: "tenant-123"
  APIKey: "01234567890123456789012345678901"
  Endpoints:
    Metrics: "http://localhost:8080/metrics"
    HealthCheck: "http://localhost:8080/health"
    KeyValidation: "http://localhost:8080/validate"
LogFilePath: "./agent.log"
LogSettings:
  Level: "info"
  Format: "json"
  MaxSize: 100
  MaxBackups: 3
  MaxAge: 28
  Compress: true
Interval: 30
Kafka:
  Brokers:
    - "localhost:9092"
  Topic: "system-metrics"
  TenantTopic: "tenant-{id}-metrics"
  SecurityProtocol: "plaintext"
Monitors:
  CPU: true
  Memory: true
  Disk: true
  Network: false
  Process: true
HTTP:
  Endpoint: "http://localhost:8080"
  StorageDir: "./storage"
  RetryAttempts: 3
  RetryDelay: 5
  Timeout: 30
  Headers:
    TenantID: "X-Tenant-ID"
    APIKey: "X-API-Key"
StorageDir: "./storage"
Thresholds:
  CPU: 80
  Memory: 85
  Disk: 90
  NetworkUtilization: 80
Storage:
  MaxStoragePerTenant: 1024
  RetentionPeriod: 7
  CompressOldData: true
Security:
  EncryptData: true
  ValidateSSL: true
  AllowedIPs: []
`,
			expectError: false,
			expected: Config{
				Version: "1.0.0",
				Tenant: TenantConfig{
					ID:     "tenant-123",
					APIKey: "01234567890123456789012345678901",
					Endpoints: struct {
						Metrics       string `yaml:"Metrics"`
						HealthCheck   string `yaml:"HealthCheck"`
						KeyValidation string `yaml:"KeyValidation"`
					}{
						Metrics:       "http://localhost:8080/metrics",
						HealthCheck:   "http://localhost:8080/health",
						KeyValidation: "http://localhost:8080/validate",
					},
				},
				LogFilePath: "./agent.log",
				LogSettings: LogSettings{
					Level:      "info",
					Format:     "json",
					MaxSize:    100,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
				},
				Interval: 30,
				Kafka: KafkaConfig{
					Brokers:          []string{"localhost:9092"},
					Topic:            "system-metrics",
					TenantTopic:      "tenant-{id}-metrics",
					SecurityProtocol: "plaintext",
				},
				Monitors: struct {
					CPU     bool `yaml:"CPU"`
					Memory  bool `yaml:"Memory"`
					Disk    bool `yaml:"Disk"`
					Network bool `yaml:"Network"`
					Process bool `yaml:"Process"`
				}{
					CPU:     true,
					Memory:  true,
					Disk:    true,
					Network: false,
					Process: true,
				},
				HTTP: HTTPConfig{
					Endpoint:      "http://localhost:8080",
					StorageDir:    "./storage",
					RetryAttempts: 3,
					RetryDelay:    5,
					Timeout:       30,
					Headers: struct {
						TenantID string `yaml:"TenantID"`
						APIKey   string `yaml:"APIKey"`
					}{
						TenantID: "X-Tenant-ID",
						APIKey:   "X-API-Key",
					},
				},
				StorageDir: "./storage",
				Thresholds: struct {
					CPU                int `yaml:"CPU"`
					Memory             int `yaml:"Memory"`
					Disk               int `yaml:"Disk"`
					NetworkUtilization int `yaml:"NetworkUtilization"`
				}{
					CPU:                80,
					Memory:             85,
					Disk:               90,
					NetworkUtilization: 80,
				},
				Storage: StorageConfig{
					MaxStoragePerTenant: 1024,
					RetentionPeriod:     7,
					CompressOldData:     true,
				},
				Security: SecurityConfig{
					EncryptData: true,
					ValidateSSL: true,
					AllowedIPs:  []string{},
				},
			},
		},
		{
			name:       "Default config file",
			configFile: "config.yaml",
			configData: `
LogFilePath: "./agent.log"
Interval: 60
Monitors:
  CPU: true
  Memory: true
  Disk: true
  Network: true
  Process: false
HTTP:
  Endpoint: "http://localhost:8080"
  StorageDir: "./storage"
StorageDir: "./storage"
`,
			expectError: false,
			expected: Config{
				Version:     "1.0.0",
				LogFilePath: "./agent.log",
				LogSettings: LogSettings{
					Level:      "info",
					Format:     "json",
					MaxSize:    100,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
				},
				Interval: 60,
				Monitors: struct {
					CPU     bool `yaml:"CPU"`
					Memory  bool `yaml:"Memory"`
					Disk    bool `yaml:"Disk"`
					Network bool `yaml:"Network"`
					Process bool `yaml:"Process"`
				}{
					CPU:     true,
					Memory:  true,
					Disk:    true,
					Network: true,
					Process: false,
				},
				HTTP: HTTPConfig{
					Endpoint:   "http://localhost:8080",
					StorageDir: "./storage",
					Headers: struct {
						TenantID string `yaml:"TenantID"`
						APIKey   string `yaml:"APIKey"`
					}{},
				},
				StorageDir: "./storage",
				Storage: StorageConfig{
					MaxStoragePerTenant: 1024,
					RetentionPeriod:     7,
					CompressOldData:     true,
				},
				Security: SecurityConfig{
					EncryptData: true,
					ValidateSSL: true,
					AllowedIPs:  []string{},
				},
			},
		},
		{
			name:        "Missing config file",
			configFile:  "missing_config.yaml",
			configData:  "",
			expectError: false,
			expected: Config{
				Version:     "1.0.0",
				LogFilePath: "./agent.log",
				LogSettings: LogSettings{
					Level:      "info",
					Format:     "json",
					MaxSize:    100,
					MaxBackups: 3,
					MaxAge:     28,
					Compress:   true,
				},
				Interval: 60,
				Monitors: struct {
					CPU     bool `yaml:"CPU"`
					Memory  bool `yaml:"Memory"`
					Disk    bool `yaml:"Disk"`
					Network bool `yaml:"Network"`
					Process bool `yaml:"Process"`
				}{
					CPU:     true,
					Memory:  true,
					Disk:    true,
					Network: true,
					Process: false,
				},
				HTTP: HTTPConfig{
					Headers: struct {
						TenantID string `yaml:"TenantID"`
						APIKey   string `yaml:"APIKey"`
					}{},
				},
				Storage: StorageConfig{
					MaxStoragePerTenant: 1024,
					RetentionPeriod:     7,
					CompressOldData:     true,
				},
				Security: SecurityConfig{
					EncryptData: true,
					ValidateSSL: true,
					AllowedIPs:  []string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp config directory
			configDir := filepath.Join(t.TempDir(), "configs")
			err := os.MkdirAll(configDir, 0755)
			assert.NoError(t, err)

			// Write config file if needed
			if tt.configData != "" {
				configPath := filepath.Join(configDir, tt.configFile)
				err = os.WriteFile(configPath, []byte(tt.configData), 0644)
				assert.NoError(t, err)
			}

			// Set config path and load config
			viper.Reset()
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
			viper.AddConfigPath(configDir)

			cfg, err := LoadConfig()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Version, cfg.Version)
				if tt.expected.Tenant.ID != "" {
					assert.Equal(t, tt.expected.Tenant, cfg.Tenant)
				}
				assert.Equal(t, tt.expected.LogFilePath, cfg.LogFilePath)
				assert.Equal(t, tt.expected.LogSettings, cfg.LogSettings)
				assert.Equal(t, tt.expected.Interval, cfg.Interval)
				if tt.expected.Kafka.Topic != "" {
					assert.Equal(t, tt.expected.Kafka, cfg.Kafka)
				}
				assert.Equal(t, tt.expected.Monitors, cfg.Monitors)
				assert.Equal(t, tt.expected.HTTP, cfg.HTTP)
				assert.Equal(t, tt.expected.StorageDir, cfg.StorageDir)
				if tt.expected.Thresholds.CPU != 0 {
					assert.Equal(t, tt.expected.Thresholds, cfg.Thresholds)
				}
				assert.Equal(t, tt.expected.Storage, cfg.Storage)
				assert.Equal(t, tt.expected.Security, cfg.Security)
			}
		})
	}
}
