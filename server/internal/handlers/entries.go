package handlers

import (
	"context"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type EntriesHandler struct {
	requirementEntries db.RequirementEntries
}

func InitEntriesHandler(
	requirementEntries db.RequirementEntries,
) (*EntriesHandler, error) {
	return &EntriesHandler{
		requirementEntries: requirementEntries,
	}, nil
}

// CreateRequirementEntry godoc
// @Summary Create a new entry
func (h *EntriesHandler) AddRequirementEntry(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	h.requirementEntries.UpsertTx()

}

// This one takes IDs and two dates in query

// CreateTask godoc
// @Summary Create a new entry
func (h *EntriesHandler) GetEntries(c *gin.Context) {
}
