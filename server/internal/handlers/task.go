package handlers

import (
	"net/http"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
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

func (h *TaskHandler) Add(c *gin.Context) {

	logger.Log.Debug().Msg("Trying to add task...")

	var req dto.TaskCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Проверка обязательного поля Requirement
	// TODO: Сделать нормальную ошибку
	if req.Task.Requirement == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requirement is required"})
		return
	}

	task := models.Task{
		Name:        req.Task.Name,
		Description: req.Task.Description,
		Requirement: &models.TaskRequirement{ // Важно: инициализируем структуру!
			Type:  req.Task.Requirement.Type,
			Value: req.Task.Requirement.Value,
		},
	}

	logger.Log.Debug().Str("name", task.Name).Msg("Creating task...")

	// Get claims
	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		c.AbortWithStatusJSON(500, dto.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "Internal error",
		})
		return
	}

	claims, ok := rawClaims.(*security.CustomClaims)
	if !ok {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		c.AbortWithStatusJSON(500, dto.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "Internal error",
		})
		return
	}

	task.OwnerID = claims.UserID

	h.taskRepo.CreateTask(&task)
}
