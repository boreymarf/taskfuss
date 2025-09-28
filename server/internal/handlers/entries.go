package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		})
	}

	api.Accepted(c, res)
}

// GetRequirementEntries godoc
// @Summary Get all entries with filtering options
// @Description Get all requirement entries with filters for archived/active status and date ranges
// @Tags entries
// @Accept json
// @Produce json
// @Param archived query bool false "Filter by archived status (true/false)"
// @Param start_date query string false "Start date for date range filter (YYYY-MM-DD)"
// @Param end_date query string false "End date for date range filter (YYYY-MM-DD)"
// @Success 200 {object} api.Response{data=[]dto.RequirementEntryResponse}
// @Failure 400 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /entries [get]
func (h *EntriesHandler) GetRequirementEntries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var queryParams struct {
		Archived  *bool   `form:"archived"`
		StartDate *string `form:"start_date"`
		EndDate   *string `form:"end_date"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		api.BadRequest.SendAndAbort(c)
		return
	}

	var startDate, endDate *time.Time

	if queryParams.StartDate != nil && *queryParams.StartDate != "" {
		parsedStart, err := time.Parse("2006-01-02", *queryParams.StartDate)
		if err != nil {
			api.BadRequest.SendWithDetailsAndAbort(c, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
		startDate = &parsedStart
	} else {
		today := time.Now()
		startDate = &today
	}

	if queryParams.EndDate != nil && *queryParams.EndDate != "" {
		parsedEnd, err := time.Parse("2006-01-02", *queryParams.EndDate)
		if err != nil {
			api.BadRequest.SendWithDetailsAndAbort(c, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
		endOfDay := time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 0, parsedEnd.Location())
		endDate = &endOfDay
	} else {
		today := time.Now()
		endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, today.Location())
		endDate = &endOfDay
	}

	if startDate.After(*endDate) {
		api.BadRequest.SendWithDetailsAndAbort(c, "start_date cannot be after end_date")
		return
	}

	showArchived := false
	if queryParams.Archived != nil {
		showArchived = *queryParams.Archived
	}

	claims := security.GetClaimsFromContext(c)

	filter := service.GetRequirementEntriesQueryParams{
		ShowArchived: showArchived,
		// TODO:
		// StartDate:    startDate,
		// EndDate:      endDate,
	}

	// Get entries from service
	entries, err := h.entriesService.GetRequirementEntries(ctx, filter, claims.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Service failed to get requirement entries!")
		api.InternalServerError.SendAndAbort(c)
		return
	}

	// Convert to response DTO
	res := make([]dto.RequirementEntryResponse, 0, len(entries))

	for _, entry := range entries {
		res = append(res, dto.RequirementEntryResponse{
			ID:            entry.ID,
			RevisionUUID:  entry.RevisionUUID,
			RequirementID: entry.RequirementID,
			Date:          entry.EntryDate,
			Value:         entry.Value,
			// Add any additional fields you might have
		})
	}

	api.Success(c, res)
}

// GetRequirementEntryByID godoc
// @Summary Get specific entry by ID
// @Description Get a specific requirement entry by its ID
// @Tags entries
// @Accept json
// @Produce json
// @Param entry_id path int true "Entry ID"
// @Success 200 {object} api.Response{data=dto.RequirementEntryResponse}
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /entries/{entry_id} [get]
func (h *EntriesHandler) GetRequirementEntryByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Extract entry_id from the path parameter
	entryIDParam := c.Param("entry_id")
	entryID, err := strconv.ParseInt(entryIDParam, 10, 64)
	if err != nil {
		api.BadRequest.SendAndAbort(c)
		return
	}

	// User auth
	claims := security.GetClaimsFromContext(c)

	// Get entry from service
	entry, err := h.entriesService.GetRequirementEntryByID(ctx, entryID, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "not found") {
			api.NotFound.SendWithDetailsAndAbort(c, "Entry not found")
			return
		}
		logger.Log.Error().Err(err).Msg("Service failed to get requirement entry by ID!")
		api.InternalServerError.SendAndAbort(c)
		return
	}

	// Convert to response DTO
	res := dto.RequirementEntryResponse{
		ID:            entry.ID,
		RevisionUUID:  entry.RevisionUUID,
		RequirementID: entry.RequirementID,
		Date:          entry.EntryDate,
		Value:         entry.Value,
	}

	api.Success(c, res)
}
