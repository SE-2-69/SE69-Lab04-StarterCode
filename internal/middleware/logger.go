package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// log every request: method, path, status, time used
// (never log headers or bodies: they contain passwords, tokens, and cookies)
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		log.Printf("%-6s %-28s %d %s (%s)", c.Request.Method, c.Request.URL.Path, status, http.StatusText(status), time.Since(start).Round(time.Microsecond))
	}
}
