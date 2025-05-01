package handlers

import (
	"github.com/gin-gonic/gin"
)

// Обработчик для /api/test
func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, we can hear you.",
	})
}
