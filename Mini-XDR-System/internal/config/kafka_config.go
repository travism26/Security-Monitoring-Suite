package config

// KafkaConfig holds configuration for Kafka consumer
type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	ConsumerGroup   string   `yaml:"consumer_group"`
	Topics          []string `yaml:"topics"`
	SecurityEnabled bool     `yaml:"security_enabled"`
	SASL            struct {
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Mechanism string `yaml:"mechanism"`
	} `yaml:"sasl"`
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CAFile   string `yaml:"ca_file"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
}

// Config represents the complete application configuration
type Config struct {
	Kafka    KafkaConfig `yaml:"kafka"`
	LogLevel string      `yaml:"log_level"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Kafka: KafkaConfig{
			Brokers:       []string{"localhost:9092"},
			ConsumerGroup: "mini-xdr-consumer",
			Topics:        []string{"security-events"},
		},
		LogLevel: "info",
	}
}
