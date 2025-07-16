package security

import (
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/gin-gonic/gin"
)

func GetClaimsFromContext(c *gin.Context) *CustomClaims {
	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims")
		c.AbortWithStatusJSON(500, dto.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "Internal error",
		})
		return nil
	}

	claims, ok := rawClaims.(*CustomClaims)
	if !ok {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		c.AbortWithStatusJSON(500, dto.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "Internal error",
		})
		return nil
	}

	return claims
}
