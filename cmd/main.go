package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"task-management-api/config"
	"task-management-api/internal/api"
	"task-management-api/internal/api/handlers"
	"task-management-api/internal/repository"
	"task-management-api/internal/service"
	"task-management-api/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewMariaDBConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(taskService)
	userHandler := handlers.NewUserHandler(userService)

	// Set up Gin router
	router := gin.Default()

	// Set up routes
	api.SetupRoutes(router, taskHandler, userHandler)

	// Start the server
	log.Printf("Starting server on %s", cfg.Server.Addr)
	if err := router.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}