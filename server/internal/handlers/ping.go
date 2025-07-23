package handlers

import (
	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	api.Success(c, gin.H{"message": "pong"})
}
