package routes

import (
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/test", handlers.TestHandler) // GET /api/test
	}
	logger.Log.Info().Msg("API end points are connected.")
}
