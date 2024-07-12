package service

import (
	"task-management-api/internal/models"
	"task-management-api/internal/repository"
)

type TaskService interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id int) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(task *models.Task) error {
	return s.repo.CreateTask(task)
}

func (s *taskService) GetTaskByID(id int) (*models.Task, error) {
	return s.repo.GetTaskByID(id)
}

func (s *taskService) GetAllTasks() ([]*models.Task, error) {
	return s.repo.GetAllTasks()
}

func (s *taskService) UpdateTask(task *models.Task) error {
	return s.repo.UpdateTask(task)
}

func (s *taskService) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}