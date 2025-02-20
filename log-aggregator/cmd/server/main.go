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

	_ "github.com/lib/pq"

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

	// Set database connection for postgres package
	postgres.SetDB(db)

	// Initialize repositories
	logRepo := postgres.NewLogRepository(db)
	logRepo.SetBatchSize(cfg.Database.BatchSize)
	processRepo := postgres.NewProcessRepository(db)
	alertRepo := postgres.NewAlertRepository(db, cfg.Features.MultiTenancy.Enabled)

	log.Printf("Log repository configured with batch size: %d", cfg.Database.BatchSize)

	// Initialize services
	logService := service.NewLogService(logRepo, service.LogServiceConfig{
		Environment: cfg.LogService.Environment,
		Application: cfg.LogService.Application,
		Component:   cfg.LogService.Component,
		Cache: struct {
			Enabled      bool
			TTL          time.Duration
			TimeRangeTTL time.Duration
		}{
			Enabled:      cfg.Cache.Enabled,
			TTL:          time.Duration(cfg.Cache.TTL) * time.Minute,
			TimeRangeTTL: time.Duration(cfg.Cache.TimeRangeTTL) * time.Minute,
		},
	})
	alertService := service.NewAlertService(alertRepo, &service.AlertServiceConfig{
		OrganizationID: cfg.Organization.ID,
		SystemMemory:   16 * 1024 * 1024 * 1024, // 16GB default
		TimeNowFn:      time.Now,
	})

	// Start Kafka consumer
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		logService,
		alertService,
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

	// Initialize HTTP server with minimal middleware
	router := gin.New()
	router.Use(gin.Recovery()) // Add recovery middleware globally for safety

	// Register health check endpoints directly on the router (no auth required)
	handler.RegisterHealthRoutes(router)

	// Create API router group with full middleware stack
	apiRouter := router.Group("/api/v1")
	apiRouter.Use(
		middleware.CORS(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Tenant(), // Add tenant middleware for multi-tenancy support
	)

	// Register API routes on the authenticated router group
	logHandler := handler.NewLogHandler(logService)
	alertHandler := handler.NewAlertHandler(alertService)

	// Register routes without the /api/v1 prefix since it's already in the group
	logs := apiRouter.Group("/logs")
	{
		logs.GET("", logHandler.ListLogs)
		logs.GET("/:id", logHandler.GetLog)
		logs.GET("/range", logHandler.ListLogsByTimeRange)
		logs.POST("", logHandler.StoreLog)
		logs.POST("/batch", logHandler.StoreBatchLogs)
	}

	alerts := apiRouter.Group("/alerts")
	{
		alerts.GET("", alertHandler.ListAlerts)
		alerts.GET("/:id", alertHandler.GetAlert)
		alerts.GET("/status/:status", alertHandler.ListAlertsByStatus)
		alerts.GET("/severity/:severity", alertHandler.ListAlertsBySeverity)
		alerts.GET("/trends", alertHandler.GetAlertTrends)
		alerts.PUT("/:id/status", alertHandler.UpdateAlertStatus)
	}

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

	// Set connection pool settings from configuration
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	// Log connection pool settings
	log.Printf("Database connection pool configured with: MaxOpenConns=%d, MaxIdleConns=%d, ConnMaxLifetime=%dm",
		cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetime)

	return db, nil
}
