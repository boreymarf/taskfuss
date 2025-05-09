package handlers

import "github.com/gin-gonic/gin"

func UserInfoHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello, we can hear you.",
	})
}
