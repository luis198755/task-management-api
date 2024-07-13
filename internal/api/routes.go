package api

import (
	"github.com/gin-gonic/gin"
	"task-management-api/internal/api/handlers"
	"task-management-api/internal/api/middleware"
)

func SetupRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler, userHandler *handlers.UserHandler) {
	v1 := router.Group("/api/v1")
	{
		// Public routes
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.RegisterUser)
			users.POST("/login", userHandler.Login)
		}

		// Protected routes
		authenticated := v1.Group("/")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// User routes
			users := authenticated.Group("/users")
			{
				users.GET("", userHandler.ListUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Task routes
			tasks := authenticated.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetAllTasks)
				tasks.GET("/:id", taskHandler.GetTaskByID)
				tasks.POST("", taskHandler.CreateTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
			}
		}
	}
}