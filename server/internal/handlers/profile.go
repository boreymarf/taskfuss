package handlers

import (
	"context"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userRepo db.Users
}

func InitProfileHandler(userRepo db.Users) (*ProfileHandler, error) {
	return &ProfileHandler{userRepo: userRepo}, nil
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieves the authenticated user's profile information
// @Tags profile
// @Security ApiKeyAuth
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dto.ProfileResponse "Profile retrieved successfully"
// @Failure 401 {object} api.Error "Unauthorized (code: UNAUTHORIZED)"
// @Failure 404 {object} api.Error "Profile not found (code: PROFILE_NOT_FOUND)"
// @Failure 500 {object} api.Error "Internal server error (code: INTERNAL_ERROR)"
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		api.InternalServerError.SendAndAbort(c)
		return
	}

	claims, ok := rawClaims.(*security.CustomClaims)
	if !ok {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		api.InternalServerError.SendAndAbort(c)
		return
	}

	uc, err := h.userRepo.GetContextByID(ctx, claims.UserID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to get user context in the profile handler")
		api.InternalServerError.SendAndAbort(c)
	}

	modelsUser, err := h.userRepo.Get(ctx, uc).WithIDs(claims.UserID).First()

	dtoUser := dto.User{
		Id:        modelsUser.ID,
		Username:  modelsUser.Username,
		CreatedAt: modelsUser.CreatedAt,
	}

	dtoProfileResponse := dto.ProfileResponse{
		User: dtoUser,
	}

	api.Success(c, dtoProfileResponse)
}
