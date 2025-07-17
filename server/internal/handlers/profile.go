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

func (h *ProfileHandler) GetProfile(c *gin.Context) {

	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		api.InternalServerError.Send(c)
		return
	}

	claims, ok := rawClaims.(*security.CustomClaims)
	if !ok {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		api.InternalServerError.Send(c)
		return
	}

	var user models.User
	h.userRepo.GetUserByID(claims.UserID, &user)

	c.JSON(200, dto.User{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}
