package handlers

import (
	"net/http"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/boreymarf/task-fuss/server/internal/service"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	userRepo             *db.UserRepository
	taskRepo             *db.TaskRepository
	taskEntryRepo        *db.TaskEntryRepository
	requirementRepo      *db.RequirementRepository
	requirementEntryRepo *db.RequirementEntryRepository
	taskService          *service.TaskService
}

func InitTaskHandler(
	userRepo *db.UserRepository,
	taskRepo *db.TaskRepository,
	taskEntryRepo *db.TaskEntryRepository,
	requirementRepo *db.RequirementRepository,
	requirementEntryRepo *db.RequirementEntryRepository,
	taskService *service.TaskService,

) (*TaskHandler, error) {
	return &TaskHandler{
		userRepo:        userRepo,
		taskRepo:        taskRepo,
		taskEntryRepo:   taskEntryRepo,
		requirementRepo: requirementRepo,
		taskService:     taskService,
	}, nil
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	// c.JSON(200, gin.H{
	// 	"message": "Hello, we can hear you.",
	// })
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {

	idParam := c.Param("id")
	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Get claims
	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims")
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

	var task models.Task
	err = h.taskRepo.GetTaskByID(taskID, &task)
	if err != nil {
		// TODO: This needs to be more specific
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if task.OwnerID != claims.UserID {
		// TODO: This needs to be more specific
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

}

func (h *TaskHandler) Add(c *gin.Context) {

	var req dto.TaskAddRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Get claims
	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims")
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
	if err := h.taskService.AddTask(&req, claims.UserID); err != nil {

		// TODO: Make normal errors
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// task := models.Task{
	// 	Name:        req.Task.Name,
	// 	Description: req.Task.Description,
	// 	Requirement: &models.TaskRequirement{ // Важно: инициализируем структуру!
	// 		Type:  req.Task.Requirement.Type,
	// 		Value: req.Task.Requirement.Value,
	// 	},
	// }
	//
	// logger.Log.Debug().Str("name", task.Name).Msg("Creating task...")
	//
	// // Get claims
	// rawClaims, exists := c.Get("userClaims")
	// if !exists {
	// 	logger.Log.Error().Msg("Failed to get claims in the profile handler")
	// 	c.AbortWithStatusJSON(500, dto.GenericError{
	// 		Code:    "INTERNAL_ERROR",
	// 		Message: "Internal error",
	// 	})
	// 	return
	// }
	//
	// claims, ok := rawClaims.(*security.CustomClaims)
	// if !ok {
	// 	logger.Log.Error().Msg("Failed to get claims in the profile handler")
	// 	c.AbortWithStatusJSON(500, dto.GenericError{
	// 		Code:    "INTERNAL_ERROR",
	// 		Message: "Internal error",
	// 	})
	// 	return
	// }
	//
	// task.OwnerID = claims.UserID
	//
	// if err := h.taskRepo.CreateTask(&task); err != nil {
	// 	logger.Log.Error().Err(err).Msg("Failed to handle task add!")
	// 	c.AbortWithStatusJSON(500, dto.GenericError{
	// 		Code:    "INTERNAL_ERROR",
	// 		Message: "Internal error",
	// 	})
	// 	return
	// }
}
