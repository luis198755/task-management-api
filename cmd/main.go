package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"task-management-api/config"
	"task-management-api/internal/api"
	"task-management-api/internal/api/handlers"
	"task-management-api/internal/repository"
	"task-management-api/internal/service"
	"task-management-api/pkg/database"
	"task-management-api/internal/api/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up logging
	setupLogging(cfg.Log)

	// Initialize database connection
	db, err := database.NewMariaDBConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Initialize repository
	taskRepo := repository.NewTaskRepository(db)

	// Initialize service
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(taskService)

	// Set up Gin router
	router := gin.Default()
	router.Use(middleware.Logger())

	// Set up CORS
	corsConfig := cors.DefaultConfig()
	if len(cfg.CORS.AllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.CORS.AllowedOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}
	corsConfig.AllowMethods = cfg.CORS.AllowedMethods
	corsConfig.AllowHeaders = cfg.CORS.AllowedHeaders
	corsConfig.MaxAge = cfg.CORS.MaxAge
	router.Use(cors.New(corsConfig))

	// Set up routes
	api.SetupRoutes(router, taskHandler)

	// Create a new http.Server
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupLogging(cfg config.LogConfig) {
	// Set up logging based on configuration
	// This is a simple setup, you might want to use a more robust logging library
	if cfg.Format == "json" {
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
	}

	// Set log level (this is a simplified version, you might want to implement more sophisticated logging)
	switch cfg.Level {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "info":
		gin.SetMode(gin.ReleaseMode)
	case "error":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}