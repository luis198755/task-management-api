package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-management-api/internal/models"
	"task-management-api/internal/service"
	"regexp"
	"github.com/go-playground/validator/v10"
	"strings"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterUser handles user registration
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var newUser models.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": getValidationErrors(err)})
		return
	}

	// Check for missing required fields
	errors := make(map[string]string)
	if newUser.Username == "" {
		errors["username"] = "required"
	}
	if newUser.Email == "" {
		errors["email"] = "required"
	}
	if newUser.Password == "" {
		errors["password"] = "required"
	}
	if newUser.Role == "" {
		errors["role"] = "required"
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": errors})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": gin.H{"email": "invalid email format"}})
		return
	}

	// Validate password length
	if len(newUser.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": gin.H{"password": "too short, minimum 8 characters required"}})
		return
	}

	user, err := h.userService.CreateUser(&newUser)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

func getValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[strings.ToLower(e.Field())] = e.Tag()
		}
	} else {
		errors["general"] = err.Error()
	}
	return errors
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	var credentials models.UserCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, token, err := h.userService.Authenticate(&credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": token})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser updates a user's information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updates models.UpdateUser
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.userService.UpdateUser(id, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ListUsers retrieves a list of users
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	users, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

