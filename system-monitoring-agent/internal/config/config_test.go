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
			name:       "Valid config file",
			configFile: "config.yaml",
			configData: `
LogFilePath: /var/log/agent.log
Interval: 30
Monitors:
  CPU: true
  Memory: true
  Disk: true
  Network: false
HTTP:
  Endpoint: http://localhost:8080
StorageDir: /data/metrics
`,
			expectError: false,
			expected: Config{
				LogFilePath: "/var/log/agent.log",
				Interval:    30,
				Monitors: struct {
					CPU     bool `yaml:"CPU"`
					Memory  bool `yaml:"Memory"`
					Disk    bool `yaml:"Disk"`
					Network bool `yaml:"Network"`
				}{
					CPU:     true,
					Memory:  true,
					Disk:    true,
					Network: false,
				},
				HTTP: HTTPConfig{
					Endpoint: "http://localhost:8080",
				},
				StorageDir: "/data/metrics",
			},
		},
		{
			name:        "Missing config file - use defaults",
			configFile:  "missing_config.yaml",
			configData:  "",
			expectError: false,
			expected: Config{
				LogFilePath: "./agent.log",
				Interval:    60,
				Monitors: struct {
					CPU     bool `yaml:"CPU"`
					Memory  bool `yaml:"Memory"`
					Disk    bool `yaml:"Disk"`
					Network bool `yaml:"Network"`
				}{
					CPU:     true,
					Memory:  true,
					Disk:    false,
					Network: false,
				},
				StorageDir: "./metrics_data",
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
				assert.Equal(t, tt.expected, *cfg)
			}
		})
	}
}
