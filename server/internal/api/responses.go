package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data any) {
	sendResponse(c, 200, data)
}

func Created(c *gin.Context, data any) {
	sendResponse(c, 201, data)
}

func Accepted(c *gin.Context, data any) {
	sendResponse(c, 202, data)
}

// Unified response handler
func sendResponse(c *gin.Context, status int, data any) {
	// Add timing headers if available
	if start, exists := c.Get("request_start"); exists {
		if startTime, ok := start.(time.Time); ok {
			c.Header("Request-Latency", time.Since(startTime).String())
		}
	}

	c.JSON(status, data)
}
