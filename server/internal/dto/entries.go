package dto

import "time"

type RequirementEntryRequest struct {
	RequirementID int64     `json:"requirement_id" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	Value         string    `json:"value" binding:"required"`
}
