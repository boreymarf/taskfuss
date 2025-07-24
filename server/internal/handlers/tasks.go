package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
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

func (h *TaskHandler) CreateTask(c *gin.Context) {

	var req dto.TaskCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	claims := security.GetClaimsFromContext(c)

	createdTask, err := h.taskService.CreateTask(&req, claims.UserID)
	if err != nil {
		api.InternalServerError.SendAndAbort(c)
	}

	api.Accepted(c, createdTask)
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

	api.Success(c, dto.GetAllTasksResponse{
		Tasks: tasks,
	})
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {

	idParam := c.Param("task_id")

	taskId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Log.Warn().Str("task_id", idParam).Msg("Tried to parse bad task id")
		api.BadRequest.SendAndAbort(c)
	}

	claims := security.GetClaimsFromContext(c)

	dtoTask, err := h.taskService.GetTaskByID(taskId, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.NotFound.SendWithDetailsAndAbort(c, gin.H{"error": "Task not found"})
		} else {
			api.InternalServerError.SendAndAbort(c)
		}
		return
	}

	data := dto.GetTaskByIDResponse{
		Task: dtoTask,
	}

	api.Success(c, data)
}

// func (h *TaskHandler) DeleteTaskByID(c *gin.Context) {
//
// 	idParam := c.Param("task_id")
//
// 	taskId, err := strconv.ParseInt(idParam, 10, 64)
// 	if err != nil {
// 		logger.Log.Warn().Str("task_id", idParam).Msg("Tried to parse bad task id")
// 		api.BadRequest.SendAndAbort(c)
// 	}
//
// 	claims := security.GetClaimsFromContext(c)
//
// 	dtoTask, err := h.taskService.GetTaskByID(taskId, claims.UserID)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			api.NotFound.SendWithDetailsAndAbort(c, gin.H{"error": "Task not found"})
// 		} else {
// 			api.InternalServerError.SendAndAbort(c)
// 		}
// 		return
// 	}
//
// 	api.Success(c, d)
//
// }
