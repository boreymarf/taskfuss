package main

import (
	"fmt"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	logger.Log.Info().Msg("Starting server...")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки .env")
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	routes.SetupAPIRoutes(r)

	// Запуск на порту 4000
	if err := r.Run(":4000"); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to start server")
	}
}
