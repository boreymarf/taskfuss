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

	nodes, err := h.entriesService.UpsertRequirementEntry(ctx, modelsEntry, claims.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Service failed to upsert a requirement entry!")
		api.InternalServerError.SendAndAbort(c)
	}

	// I decided to just return a slice of entries instead of a tree lmao
	res := make([]dto.RequirementEntryResponse, 0, len(nodes))

	for _, node := range nodes {
		content := node.Content.(*models.RequirementEntry)
		res = append(res, dto.RequirementEntryResponse{
			ID:            content.ID,
			RevisionUUID:  content.RevisionUUID,
			RequirementID: content.RequirementID,
			Date:          content.EntryDate,
			Value:         content.Value,
			Children:      []dto.RequirementEntryResponse{},
		})
	}

	api.Accepted(c, res)
}

// This one takes IDs and two dates in query

// CreateTask godoc
// @Summary Create a new entry
func (h *EntriesHandler) GetEntries(c *gin.Context) {
}
