package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rickjms/security-monitoring-suite/mini-xdr/internal/config"
	"github.com/rickjms/security-monitoring-suite/mini-xdr/internal/kafka"
	"github.com/rickjms/security-monitoring-suite/mini-xdr/internal/service"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create event processors
	systemMetricsProcessor := service.NewSystemMetricsProcessor(logger)
	networkTrafficProcessor := service.NewNetworkTrafficProcessor(logger)

	// Create event service
	eventService := service.NewEventService(logger, systemMetricsProcessor, networkTrafficProcessor)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.ConsumerGroup,
		cfg.Kafka.Topics,
		eventService,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to create Kafka consumer", zap.Error(err))
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
	}()

	// Start consuming events
	logger.Info("Starting Mini-XDR server",
		zap.Strings("brokers", cfg.Kafka.Brokers),
		zap.Strings("topics", cfg.Kafka.Topics))

	if err := consumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start consumer", zap.Error(err))
	}

	<-ctx.Done()
	logger.Info("Shutting down Mini-XDR server")

	if err := consumer.Close(); err != nil {
		logger.Error("Error closing consumer", zap.Error(err))
	}
}

func loadConfig(path string) (*config.Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := config.DefaultConfig()
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLogger(level string) (*zap.Logger, error) {
	var cfg zap.Config

	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	return cfg.Build()
}
