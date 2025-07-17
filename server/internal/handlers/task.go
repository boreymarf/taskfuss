package handlers

import (
	"net/http"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
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
	// c.JSON(200, gin.H{ "message": "Hello, we can hear you.",
	// })
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {

	idParam := c.Param("id")

	taskID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	claims := security.GetClaimsFromContext(c)

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

func (h *TaskHandler) CreateTask(c *gin.Context) {

	var req dto.TaskAddRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	claims := security.GetClaimsFromContext(c)

	// Get claims
	if err := h.taskService.AddTask(&req, claims.UserID); err != nil {

		// TODO: Make normal errors
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
}

func (h *TaskHandler) GetRequirements(c *gin.Context) {

}
