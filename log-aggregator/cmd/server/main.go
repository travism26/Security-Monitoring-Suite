package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/config"
	"github.com/travism26/log-aggregator/internal/handler"
	"github.com/travism26/log-aggregator/internal/kafka"
	"github.com/travism26/log-aggregator/internal/middleware"
	"github.com/travism26/log-aggregator/internal/repository/postgres"
	"github.com/travism26/log-aggregator/internal/service"
)

// Add these variables at the package level, before the main function
var (
	version    = "dev"
	commitHash = "unknown"
	buildTime  = "unknown"
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
	log.Printf("Version: %s, Commit: %s, Built at: %s", version, commitHash, buildTime)

	// Create a context that will be canceled on shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := initializeDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	logRepo := postgres.NewLogRepository(db)
	processRepo := postgres.NewProcessRepository(db)

	// Initialize services
	logService := service.NewLogService(logRepo)

	// Start Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		logService,
		processRepo,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Start consumer in a goroutine with context
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
			cancel() // Cancel context on consumer error
		}
	}()

	// Initialize HTTP server with middleware
	router := gin.Default()
	router.Use(
		middleware.CORS(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
	)

	// Register health check endpoints
	handler.RegisterHealthRoutes(router)

	// Register API routes
	logHandler := handler.NewLogHandler(logService)
	handler.RegisterRoutes(router, logHandler)

	// Create HTTP server with timeout configurations
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v", err)
			cancel() // Cancel context on server error
		}
	}()

	// Wait for shutdown signal
	<-signalChan
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Cancel the context to stop the Kafka consumer
	cancel()
	log.Println("Server stopped gracefully")
}

func initializeDB(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
