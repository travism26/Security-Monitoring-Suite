package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Endpoint  string `yaml:"endpoint"`
	Interval  int    `yaml:"interval"`
	Scenarios struct {
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

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
