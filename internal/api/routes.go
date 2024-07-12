package api

import (
	"github.com/gin-gonic/gin"
	"task-management-api/internal/api/handlers"
)

func SetupRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler) {
	v1 := router.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetAllTasks)
			tasks.GET("/:id", taskHandler.GetTaskByID)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}
}