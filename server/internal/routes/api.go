package routes

import (
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	api := router.Group("/api")
	{
		api.GET("/test", handlers.TestHandler) // GET /api/test
		api.POST("/register", authHandler.RegisterHandler)
		api.POST("/login", authHandler.LoginHandler)
		api.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{"isUp": true})
		})

		auth := router.Group("")
		auth.GET("/user", middleware.Auth())
	}
	logger.Log.Info().Msg("API end points are connected!")
}
