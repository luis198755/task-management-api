package models

import "time"

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title" binding:"required,min=1,max=100"`
	Description string     `json:"description" binding:"max=500"`
	Status      TaskStatus `json:"status" binding:"required,oneof=TODO IN_PROGRESS DONE"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}