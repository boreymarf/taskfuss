package handlers

import (
	"fmt"
	"net/http"

	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userRepo *db.UserRepository
}

func InitAuthHanlder(userRepo *db.UserRepository) (*AuthHandler, error) {
	return &AuthHandler{userRepo: userRepo}, nil
}

func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest
	var errorMessage string

	if err := c.ShouldBindJSON(&req); err != nil {

		if fieldErr, ok := err.(validator.ValidationErrors); ok { // If validation fails
			firstErr := fieldErr[0] // We take first error we encounter

			switch firstErr.Tag() {
			case "required":
				errorMessage = fmt.Sprintf("%s is required", firstErr.Field())
			case "email":
				errorMessage = "Invalid email format"
			case "min":
				errorMessage = fmt.Sprintf("%s must be at least %s characters",
					firstErr.Field(), firstErr.Param())
			default:
				errorMessage = "Invalid request data"
			}
		} else {
			// Handle other JSON errors
			errorMessage = "Invalid request format"
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Hashing
	passwordHash, err := security.HashPassword(req.Password)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("name", req.Name).
			Str("email", req.Email).
			Msg("Failed to hash password for new user")

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  "INTERNAL_ERROR",
			"text":  "Internal server error",
			"error": err,
		})

		return
	}

	// Creating user
	var user models.User
	user.Name = req.Name
	user.Email = req.Email
	user.PasswordHash = passwordHash

	err = h.userRepo.CreateUser(&user)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("name", req.Name).
			Str("email", req.Email).
			Msg("Failed to create new user in the database")

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  "INTERNAL_ERROR",
			"text":  "Internal server error",
			"error": err,
		})

		return
	}

	c.JSON(200, gin.H{"result": "success"})
}
