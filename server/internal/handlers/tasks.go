package handlers

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/boreymarf/task-fuss/server/internal/service"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func InitTaskHandler(
	taskService *service.TaskService,

) (*TaskHandler, error) {
	return &TaskHandler{
		taskService: taskService,
	}, nil
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param CreateTaskRequest body dto.CreateTaskRequest true "Task creation data"
// @Success 201 {object} dto.CreateTaskResponse "Task successfully created"
// @Failure 400 {object} api.Error "Invalid request format"
// @Failure 401 {object} api.Error "Unauthorized"
// @Failure 500 {object} api.Error "Internal server error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	claims := security.GetClaimsFromContext(c)

	createdTask, err := h.taskService.CreateTask(ctx, &req, claims.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Service failed to create a new task!")
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

// GetAllTasks godoc
// @Summary Get all tasks with filtering options
// @Description Retrieves tasks based on filter criteria (active/archived/completed) and detail level
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param detailLevel query string false "Detail level" Enums(minimal, standard, full)
// @Param showActive query boolean false "Include active tasks (default: true)"
// @Param showArchived query boolean false "Include archived tasks (default: false)"
// @Param showCompleted query boolean false "Include completed tasks (default: true)"
// @Success 200 {object} dto.GetAllTasksResponse "List of tasks"
// @Failure 400 {object} api.Error "Invalid query parameters"
// @Failure 401 {object} api.Error "Unauthorized"
// @Failure 500 {object} api.Error "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var queryParams GetAllTasksQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		api.InvalidQuery.SendAndAbort(c)
	}

	opts := service.GetAllTasksQueryParams{
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

	tasks, err := h.taskService.GetAllTasks(ctx, &opts, claims.UserID)
	if err != nil {
		logger.Log.Err(err).Msg("Failed to get all tasks")
		api.InternalServerError.SendAndAbort(c)
	}

	api.Success(c, dto.GetAllTasksResponse{
		Tasks: *tasks,
	})
}

// GetTaskByID godoc
// @Summary Get a task by ID
// @Description Retrieves a single task by its unique identifier
// @Tags tasks
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param task_id path string true "Task ID" Format(uuid)
// @Success 200 {object} dto.GetTaskByIDResponse "Task details"
// @Failure 400 {object} api.Error "Invalid task ID format"
// @Failure 401 {object} api.Error "Unauthorized"
// @Failure 404 {object} api.Error "Task not found"
// @Failure 500 {object} api.Error "Internal server error"
// @Router /tasks/{task_id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	idParam := c.Param("task_id")

	taskId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Log.Warn().Str("task_id", idParam).Msg("Tried to parse bad task id")
		api.BadRequest.SendAndAbort(c)
		return
	}

	claims := security.GetClaimsFromContext(c)

	dtoTask, err := h.taskService.GetTask(ctx, taskId, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.NotFound.SendWithDetailsAndAbort(c, gin.H{"error": "Task not found"})
		} else {
			api.InternalServerError.SendAndAbort(c)
		}
		return
	}

	data := dto.GetTaskByIDResponse{
		Task: *dtoTask,
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
