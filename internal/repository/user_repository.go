package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"task-management-api/internal/models"
)

type UserRepository interface {
	CreateUser(user *models.NewUser) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id int, updates *models.UpdateUser) error
	DeleteUser(id int) error
	ListUsers(offset, limit int) ([]*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(newUser *models.NewUser) (*models.User, error) {
	query := `INSERT INTO users (username, email, password_hash, full_name, role, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	
	result, err := r.db.Exec(query, newUser.Username, newUser.Email, newUser.Password, newUser.FullName, newUser.Role)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	return r.GetUserByID(int(id))
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
			  FROM users WHERE id = ?`
	
	var user models.User
	var createdAt, updatedAt []uint8
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, 
		&user.IsActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}
	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %v", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
			  FROM users WHERE username = ?`
	
	var user models.User
	var createdAt, updatedAt []uint8
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, 
		&user.IsActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}
	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %v", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
			  FROM users WHERE email = ?`
	
	var user models.User
	var createdAt, updatedAt []uint8
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, 
		&user.IsActive, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}
	user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %v", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(id int, updates *models.UpdateUser) error {
	query := `UPDATE users SET `
	args := []interface{}{}

	if updates.Email != nil {
		query += `email = ?, `
		args = append(args, *updates.Email)
	}
	if updates.FullName != nil {
		query += `full_name = ?, `
		args = append(args, *updates.FullName)
	}
	if updates.Role != nil {
		query += `role = ?, `
		args = append(args, *updates.Role)
	}
	if updates.IsActive != nil {
		query += `is_active = ?, `
		args = append(args, *updates.IsActive)
	}

	query = query[:len(query)-2] + ` WHERE id = ?`
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) ListUsers(offset, limit int) ([]*models.User, error) {
	query := `SELECT id, username, email, password_hash, full_name, role, is_active, created_at, updated_at 
			  FROM users LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var createdAt, updatedAt []uint8
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, 
			&user.IsActive, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		
		user.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}
		user.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			return nil, fmt.Errorf("error parsing updated_at: %v", err)
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning all rows: %v", err)
	}

	return users, nil
}