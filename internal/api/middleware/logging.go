package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Printf("Request - Method: %s | Status: %d | Latency: %v | Path: %s",
			c.Request.Method,
			c.Writer.Status(),
			latency,
			c.Request.URL.Path,
		)
	}
}