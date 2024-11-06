// internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogFilePath string
	Interval    int // Polling interval in seconds
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Set default values
	viper.SetDefault("LogFilePath", "./agent.log")
	viper.SetDefault("Interval", 60)

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
