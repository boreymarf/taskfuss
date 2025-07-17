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

			protected.GET("/tasks", taskHandler.ListTasks)            // Will list all tasks with names, requrements only, end date, start date
			protected.GET("/tasks/:task_id", taskHandler.GetTaskByID) // Get other info of the task like description
			protected.PUT("/tasks/:task_id")                          // Update task
			protected.POST("/tasks", taskHandler.CreateTask)          // Create a task

			protected.GET("/requirements/entries", taskHandler.GetRequirements) // GET /requirements/entries?start=2024-01-01T00:00:00&end=2024-01-31T23:59:59
			protected.POST("/requirements/:requirement_id/entries")             // Create an entry for any requirement
			protected.PUT("/entries/:entry_id")                                 // Update entry
		}

	}
	logger.Log.Info().Msg("API end points are connected!")
}
