package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user information in the context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}