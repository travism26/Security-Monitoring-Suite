// internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogFilePath string `yaml:"LogFilePath"`
	Interval    int    `yaml:"Interval"` // Polling interval in seconds
	Monitors    struct {
		CPU     bool `yaml:"CPU"`
		Memory  bool `yaml:"Memory"`
		Disk    bool `yaml:"Disk"`
		Network bool `yaml:"Network"`
	} `yaml:"Monitors"`
	HTTP       HTTPConfig
	StorageDir string `yaml:"StorageDir"`
}

type HTTPConfig struct {
	Endpoint string `yaml:"Endpoint"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Set default values
	viper.SetDefault("LogFilePath", "./agent.log")
	viper.SetDefault("Interval", 60)
	viper.SetDefault("Monitors.CPU", true)
	viper.SetDefault("Monitors.Memory", true)
	viper.SetDefault("Monitors.Disk", false)
	viper.SetDefault("Monitors.Network", false)
	viper.SetDefault("StorageDir", "./metrics_data")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found; using defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
