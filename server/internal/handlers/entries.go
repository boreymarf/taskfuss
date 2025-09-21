package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/boreymarf/task-fuss/server/internal/service"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
)

type EntriesHandler struct {
	entriesService *service.EntriesService
}

func InitEntriesHandler(
	entriesService *service.EntriesService,
) (*EntriesHandler, error) {
	return &EntriesHandler{
		entriesService: entriesService,
	}, nil
}

// CreateRequirementEntry godoc
// @Summary Create a new entry
func (h *EntriesHandler) AddRequirementEntry(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Extract requirement_id from the path parameter
	requirementIDParam := c.Param("requirement_id")
	requirementID, err := strconv.ParseInt(requirementIDParam, 10, 64)
	if err != nil {
		api.BadRequest.SendAndAbort(c)
		return
	}

	// Request data
	var req dto.UpsertRequirementEntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// User auth
	claims := security.GetClaimsFromContext(c)

	strValue, ok := req.Value.(string)
	if !ok {
		strValue = fmt.Sprintf("%v", req.Value) // Convert any type to string
	}

	modelsEntry := &models.RequirementEntry{
		RequirementID: requirementID,
		EntryDate:     req.Date,
		Value:         strValue,
	}

	createdEntry, err := h.entriesService.UpsertRequirementEntry(ctx, modelsEntry, claims.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Service failed to upsert a requirement entry!")
		api.InternalServerError.SendAndAbort(c)
	}

	res := dto.RequirementEntryResponse{
		ID:            createdEntry.ID,
		RevisionUUID:  createdEntry.RevisionUUID,
		RequirementID: createdEntry.RequirementID,
		Date:          createdEntry.EntryDate,
		Value:         createdEntry.Value,
	}

	api.Accepted(c, res)
}

// This one takes IDs and two dates in query

// CreateTask godoc
// @Summary Create a new entry
func (h *EntriesHandler) GetEntries(c *gin.Context) {
}
