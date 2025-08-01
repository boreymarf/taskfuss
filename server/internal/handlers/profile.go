package handlers

import (
	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	userRepo *db.UserRepository
}

func InitProfileHandler(userRepo *db.UserRepository) (*ProfileHandler, error) {
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

	var user models.User
	h.userRepo.GetUserByID(claims.UserID, &user)

	// FIXME: This is incorrect for now
	c.JSON(200, dto.User{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}
