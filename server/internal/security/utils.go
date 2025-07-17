package security

import (
	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/gin-gonic/gin"
)

func GetClaimsFromContext(c *gin.Context) *CustomClaims {
	rawClaims, exists := c.Get("userClaims")
	if !exists {
		logger.Log.Error().Msg("Failed to get claims")
		api.InternalServerError.SendAndAbort(c)
		return nil
	}

	claims, ok := rawClaims.(*CustomClaims)
	if !ok {
		logger.Log.Error().Msg("Failed to get claims in the profile handler")
		api.InternalServerError.SendAndAbort(c)
		return nil
	}

	return claims
}
