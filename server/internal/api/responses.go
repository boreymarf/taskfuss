package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// @Description Standard API response wrapper
type Response struct {
	Data      any    `json:"data"`
	Timestamp string `json:"timestamp,omitempty" example:"2025-07-27T20:32:29+03:00"`
	Latency   string `json:"latency" example:"42.123µs"`
}

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
	resp := Response{
		Data: data,
	}

	// Add timing data if available
	if start, exists := c.Get("request_start"); exists {
		if startTime, ok := start.(time.Time); ok {
			resp.Timestamp = startTime.Format(time.RFC3339)
			resp.Latency = time.Since(startTime).String()
		}
	}

	c.JSON(status, resp)
}
