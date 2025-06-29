package handlers

import (
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	userRepo *db.UserRepository
	taskRepo *db.TaskRepository
}

func InitTaskHanlder(userRepo *db.UserRepository, taskRepo *db.TaskRepository) (*TaskHandler, error) {
	return &TaskHandler{userRepo: userRepo, taskRepo: taskRepo}, nil
}

func (*TaskHandler) Get(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, we can hear you.",
	})
}
