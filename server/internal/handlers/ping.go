package handlers

import (
	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/gin-gonic/gin"
)

// PingHandler godoc
// @Summary Server health check
// @Description Returns "pong" if the server is running
// @Tags service
// @Produce json
// @Success 200 {object} api.Response{data=dto.PongResponse} "Server is running"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	api.Success(c, dto.PongResponse{
		Message: "pong",
	})
}
