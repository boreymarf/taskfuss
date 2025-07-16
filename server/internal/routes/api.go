package routes

import (
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(
	router *gin.Engine,
	userRepo *db.UserRepository,
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
		protected.Use(middleware.Auth(userRepo))
		{
			protected.GET("/profile", profileHandler.GetProfile)
			protected.GET("/task", taskHandler.ListTasks)
			protected.GET("/tasks/:id", taskHandler.GetTaskByID)
			protected.POST("/task", taskHandler.CreateTask)
			protected.GET("/requirements", taskHandler.GetRequirements) // GET /requirements?start=2024-01-01T00:00:00&end=2024-01-31T23:59:59
			protected.
		}

	}
	logger.Log.Info().Msg("API end points are connected!")
}
