package handlers

import "github.com/gin-gonic/gin"

type EntriesHandler struct {
}

func InitEntriesHandler() (*EntriesHandler, error) {
	return &EntriesHandler{}, nil
}

// CreateRequirementEntry godoc
// @Summary Create a new entry
func (h *EntriesHandler) AddRequirementEntry(c *gin.Context) {
}

// This one takes IDs and two dates in query

// CreateTask godoc
// @Summary Create a new entry
func (h *EntriesHandler) GetEntries(c *gin.Context) {
}
