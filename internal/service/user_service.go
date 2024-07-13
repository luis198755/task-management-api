package service

import (
	"errors"
	"task-management-api/internal/models"
	"task-management-api/internal/repository"
	"task-management-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	CreateUser(newUser *models.NewUser) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id int, updates *models.UpdateUser) error
	DeleteUser(id int) error
	ListUsers(page, pageSize int) ([]*models.User, error)
	Authenticate(credentials *models.UserCredentials) (*models.User, string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(newUser *models.NewUser) (*models.User, error) {
	// Check if username already exists
	if _, err := s.userRepo.GetUserByUsername(newUser.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.userRepo.GetUserByEmail(newUser.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	newUser.Password = string(hashedPassword)

	// Create the user
	return s.userRepo.CreateUser(newUser)
}

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *userService) UpdateUser(id int, updates *models.UpdateUser) error {
	if updates.Email != nil {
		// Check if new email already exists
		if user, err := s.userRepo.GetUserByEmail(*updates.Email); err == nil && user.ID != id {
			return errors.New("email already exists")
		}
	}

	return s.userRepo.UpdateUser(id, updates)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.DeleteUser(id)
}

func (s *userService) ListUsers(page, pageSize int) ([]*models.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.userRepo.ListUsers(offset, pageSize)
}

func (s *userService) Authenticate(credentials *models.UserCredentials) (*models.User, string, error) {
	user, err := s.userRepo.GetUserByUsername(credentials.Username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}
