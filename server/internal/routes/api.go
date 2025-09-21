package routes

import (
	_ "github.com/boreymarf/task-fuss/server/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/handlers"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(
	router *gin.Engine,
	userRepo *db.Users,
	authHandler *handlers.AuthHandler,
	profileHandler *handlers.ProfileHandler,
	taskHandler *handlers.TaskHandler,
	entriesHandler *handlers.EntriesHandler,
) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{

		// For more info check OpenAPI docs
		api.GET("/ping", handlers.PingHandler)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.Auth(userRepo))
		{
			protected.GET("/profile", profileHandler.GetProfile)

			protected.GET("/tasks", taskHandler.GetAllTasks)
			protected.GET("/tasks/:task_id", taskHandler.GetTask) // Get other info of the task like description
			// protected.PUT("/tasks/:task_id")                          // Update task
			protected.POST("/tasks", taskHandler.CreateTask) // Create a task
			//
			// // protected.GET("/requirements/entries", taskHandler.GetRequirements) // GET /requirements/entries?start=2024-01-01T00:00:00&end=2024-01-31T23:59:59
			protected.POST("/requirements/:requirement_id/entries", entriesHandler.AddRequirementEntry) // Create an entry for any requirement
			// protected.GET("/requirements/:requirement_id/entries")  // Get all entries, make it with date start and end
			// protected.GET("/entries/:entry_id")                     // Get specific entry
			// protected.PUT("/entries/:entry_id")                     // Update entry
			// protected.DELETE("/entries/:entry_id")                  //
		}

	}
	logger.Log.Info().Msg("API end points are connected!")
}
