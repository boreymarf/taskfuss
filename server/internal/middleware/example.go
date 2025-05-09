package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Действия ДО обработки запроса
		start := time.Now()

		fmt.Printf("Начало обработки запроса: %s %s\n", c.Request.Method, c.Request.URL)

		// Передаем управление следующему обработчику в цепочке
		c.Next()

		// Действия ПОСЛЕ обработки запроса
		duration := time.Since(start)
		fmt.Printf("Завершение обработки. Время выполнения: %v\n", duration)
	}
}
