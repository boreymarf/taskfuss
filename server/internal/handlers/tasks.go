package handlers

import (
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/boreymarf/task-fuss/server/internal/service"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sanity-io/litter"
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

func (h *TaskHandler) CreateTask(c *gin.Context) {

	var req dto.TaskCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	claims := security.GetClaimsFromContext(c)

	if err := h.taskService.CreateTask(&req, claims.UserID); err != nil {
		api.InternalServerError.SendAndAbort(c)
	}
}

type GetAllTasksQuery struct {
	DetailLevel   string `form:"detail" binding:"omitempty,oneof=minimal basic full"`
	ShowActive    string `form:"active" binding:"omitempty,oneof=true false"`
	ShowArchived  string `form:"archived" binding:"omitempty,oneof=true false"`
	ShowCompleted string `form:"completed" binding:"omitempty,oneof=true false"`
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {

	var queryParams GetAllTasksQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		api.InvalidQuery.SendAndAbort(c)
	}

	opts := service.GetAllTasksOptions{
		DetailLevel: queryParams.DetailLevel,
	}

	opts.ShowActive = true
	opts.ShowArchived = false
	opts.ShowCompleted = true

	if queryParams.ShowActive != "" {
		opts.ShowActive = queryParams.ShowActive == "true"
	}
	if queryParams.ShowArchived != "" {
		opts.ShowArchived = queryParams.ShowArchived == "true"
	}
	if queryParams.ShowCompleted != "" {
		opts.ShowCompleted = queryParams.ShowCompleted == "true"
	}

	claims := security.GetClaimsFromContext(c)

	tasks, err := h.taskService.GetAllTasks(&opts, claims.UserID)
	if err != nil {
		logger.Log.Err(err).Msg("Failed to get all tasks")
		api.InternalServerError.SendAndAbort(c)
	}

	// TODO: Make this an API function later
	c.JSON(200, dto.GetAllTasksResponse{
		Tasks: tasks,
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {

	idParam := c.Param("task_id")

	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		api.BadRequest.SendWithDetailsAndAbort(c, gin.H{"error": "Invalid task ID"})
		return
	}

	claims := security.GetClaimsFromContext(c)

	var task models.Task
	err = h.taskRepo.GetTaskByID(taskID, &task)
	if err != nil {
		// TODO: This needs to be more specific
		api.BadRequest.SendWithDetailsAndAbort(c, gin.H{"error": "Invalid task ID"})
		return
	}

	if task.OwnerID != claims.UserID {
		// TODO: This needs to be more specific
		api.BadRequest.SendWithDetailsAndAbort(c, gin.H{"error": "Invalid task ID"})
		return
	}

	var modelTask models.Task
	err = h.taskRepo.GetTaskByID(taskID, &modelTask)
	if err != nil {
		api.InternalServerError.SendAndAbort(c)
	}

	var modelRequirements []models.Requirement
	modelRequirements, err = h.requirementRepo.GetRequirementsByTaskIDs([]int64{taskID})
	if err != nil {
		api.InternalServerError.SendAndAbort(c)
	}

	litter.Dump(modelRequirements)

}
