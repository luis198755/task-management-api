package models

import (
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	UserRoleUser  UserRole = "USER"
	UserRoleAdmin UserRole = "ADMIN"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username" binding:"required,min=3,max=50"`
	Email        string    `json:"email" binding:"required,email"`
	PasswordHash string    `json:"-"` // The "-" tag means this field won't be included in JSON output
	FullName     string    `json:"full_name" binding:"max=100"`
	Role         UserRole  `json:"role" binding:"required,oneof=USER ADMIN"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserCredentials represents the data needed for user authentication
type UserCredentials struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// NewUser represents the data needed to create a new user
type NewUser struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	FullName string   `json:"full_name" binding:"max=100"`
	Role     UserRole `json:"role" binding:"required,oneof=USER ADMIN"`
}

// UpdateUser represents the data that can be updated for a user
type UpdateUser struct {
	Email    *string   `json:"email" binding:"omitempty,email"`
	FullName *string   `json:"full_name" binding:"omitempty,max=100"`
	Role     *UserRole `json:"role" binding:"omitempty,oneof=USER ADMIN"`
	IsActive *bool     `json:"is_active"`
}