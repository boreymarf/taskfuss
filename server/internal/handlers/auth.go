package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/models"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/boreymarf/task-fuss/server/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *db.UserRepository
}

func InitAuthHanlder(userRepo *db.UserRepository) (*AuthHandler, error) {
	return &AuthHandler{userRepo: userRepo}, nil
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Hashing
	passwordHash, err := security.HashPassword(req.Password)

	if err != nil {

		logger.Log.Error().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Msg("Failed to hash password for new user")

		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_ERROR",
			"text": "Internal server error",
		})

		return
	}

	// Creating user
	var user models.User
	user.Username = req.Username
	user.Email = req.Email
	user.PasswordHash = passwordHash

	err = h.userRepo.CreateUser(&user)

	if err != nil {

		if errors.Is(err, apperrors.ErrDuplicate) {
			c.JSON(http.StatusConflict, dto.GenericError{
				Code:    "DUPLICATE_ENTRY",
				Message: "User already exists",
			})
			return
		}

		logger.Log.Error().
			Err(err).
			Str("name", req.Username).
			Str("email", req.Email).
			Msg("Failed to create new user in the database")

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  "INTERNAL_ERROR",
			"text":  "Internal server error",
			"error": err,
		})

		return
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := security.CreateToken(user.ID, []byte(secret), time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_ERROR",
			"text": "Internal server error",
		})
		return
	}

	c.JSON(200, dto.RegisterResponse{
		User: dto.User{
			Id:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		AuthToken: token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {

			logger.Log.Error().Err(syntaxErr).Send()

			c.JSON(http.StatusBadRequest, gin.H{
				"code":  "INVALID_JSON",
				"text":  "Invalid JSON syntax",
				"error": err,
			})
			return
		}

		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {

			logger.Log.Error().
				Str("field", typeErr.Field).
				Str("expected_type", typeErr.Type.String()).
				Msg("Type mismatch")

			c.JSON(http.StatusBadRequest, gin.H{
				"code": "TYPE_MISMATCH",
				"text": fmt.Sprintf("Field '%s' must be a %s", typeErr.Field, typeErr.Type),
			})
			return
		}
	}

	var user models.User

	err := h.userRepo.GetUserByEmail(req.Email, &user)

	logger.Log.Debug().Str("email", user.Email).Str("hash", user.PasswordHash).Send()

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: user does not exists")
			c.JSON(http.StatusUnauthorized, dto.GenericError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Invalid email or password",
			})
			return
		} else {
			logger.Log.Error().Str("email", req.Email).Err(err).Msg("Failed login attempt: internal server error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_ERROR",
				"text": "Internal server error",
			})
			return
		}
	}

	if req.Password == "" {

		logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: empty password")
		c.JSON(http.StatusUnauthorized, dto.GenericError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {

		logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: incorrect password")
		c.JSON(http.StatusUnauthorized, dto.GenericError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
		})
		return
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := security.CreateToken(user.ID, []byte(secret), time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_ERROR",
			"text": "Internal server error",
		})
		return
	}

	c.JSON(200, dto.LoginResponse{
		User: dto.User{
			Id:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		AuthToken: token,
	})
}
