package repository

import (
	"database/sql"
	"fmt"
	"time"
	"task-management-api/internal/models"
)

type TaskRepository interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id int) (*models.Task, error)
	GetAllTasks() ([]*models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id int) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	query := `INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)
	return nil
}

func (r *taskRepository) GetTaskByID(id int) (*models.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = ?`
	row := r.db.QueryRow(query, id)

	task := &models.Task{}
	var createdAt, updatedAt []uint8
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}
	
	// Parse the timestamps
	task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}
	task.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %v", err)
	}

	return task, nil
}

func (r *taskRepository) GetAllTasks() ([]*models.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		task := &models.Task{}
		var createdAt, updatedAt []uint8
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		
		// Parse the timestamps
		task.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}
		task.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing updated_at: %v", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning all rows: %v", err)
	}

	return tasks, nil
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?`
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}