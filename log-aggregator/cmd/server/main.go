package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/config"
	"github.com/travism26/log-aggregator/internal/handler"
	"github.com/travism26/log-aggregator/internal/kafka"
	"github.com/travism26/log-aggregator/internal/repository/postgres"
	"github.com/travism26/log-aggregator/internal/service"
)

// main is the application entry point. The startup flow is:
// 1. Load configuration
// 2. Initialize logger
// 3. Connect to database
// 4. Start Kafka consumer
// 5. Start HTTP server
//
// If any critical service fails to start, the application will log the error and exit
func main() {
	log.Println("Starting Log Aggregator Service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := sql.Open("postgres", "postgres connection string here")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize components
	repo := postgres.NewLogRepository(db)
	logService := service.NewLogService(repo)

	// Start Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		logService,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Initialize HTTP server
	router := gin.Default()
	logHandler := handler.NewLogHandler(logService)
	handler.RegisterRoutes(router, logHandler)

	// Start HTTP server
	log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
