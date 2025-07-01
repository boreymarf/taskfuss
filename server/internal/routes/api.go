package routes

import (
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	taskHandler *handlers.TaskHandler,
) {
	api := router.Group("/api")
	{
		api.GET("/test", handlers.TestHandler) // GET /api/test
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"isUp": true})
		})

		protected := api.Group("")
		protected.Use(middleware.Auth())
		{
			protected.GET("/profile", profileHandler.GetProfile)
			protected.GET("/tasks", profileHandler.GetProfile)
			protected.POST("/tasks", profileHandler.GetProfile)
		}

	}
	logger.Log.Info().Msg("API end points are connected!")
}
