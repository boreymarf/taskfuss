package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpsertRequirementEntryRequest struct {
	Date  time.Time `json:"date" binding:"required"`
	Value any       `json:"value" binding:"required"`
}

type RequirementEntryResponse struct {
	ID            int64     `json:"id"`
	RevisionUUID  uuid.UUID `json:"revision_uuid"`
	RequirementID int64     `json:"requirement_id" `
	Date          time.Time `json:"date" binding:"required"`
	Value         string    `json:"value" binding:"required"`
}
