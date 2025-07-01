package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	userRepo *db.UserRepository
	taskRepo *db.TaskRepository
}

func InitTaskHanlder(userRepo *db.UserRepository, taskRepo *db.TaskRepository) (*TaskHandler, error) {
	return &TaskHandler{userRepo: userRepo, taskRepo: taskRepo}, nil
}

func (h *TaskHandler) Get(c *gin.Context) {
	// c.JSON(200, gin.H{
	// 	"message": "Hello, we can hear you.",
	// })
}

func (h *TaskHandler) Create(c *gin.Context) {

	var req dto.TaskCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	h.taskRepo.CreateTask()
}
