package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-management-api/internal/models"
	"task-management-api/internal/service"
	"task-management-api/internal/errors"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("Invalid task ID"))
		return
	}

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		if apiErr, ok := err.(*errors.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr)
		} else {
			c.JSON(http.StatusInternalServerError, errors.NewInternalServerError("Failed to fetch task"))
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks", "details": err.Error()})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task data"})
		return
	}

	task.ID = id

	err = h.taskService.UpdateTask(&task)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	err = h.taskService.DeleteTask(id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}