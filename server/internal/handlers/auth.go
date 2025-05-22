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
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *db.UserRepository
}

func InitAuthHanlder(userRepo *db.UserRepository) (*AuthHandler, error) {
	return &AuthHandler{userRepo: userRepo}, nil
}

func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest

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

		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			logger.Log.Error().Err(validationErrors).Send()

			var response dto.ValidationError = dto.ValidationError{
				Code:    "VALIDATION_FAILED",
				Message: "Validation failed",
				Details: make([]dto.FieldError, 0),
			}

			for _, fieldError := range validationErrors {

				switch fieldError.Tag() {
				case "required":
					logger.Log.Info().
						Str("ip", c.ClientIP()).
						Str("field", fieldError.Field()).
						Msg("Failed registration attempt: missing required field")
					response.Details = append(response.Details, dto.FieldError{
						Field:   fieldError.Field(),
						Code:    "REQUIRED",
						Message: fmt.Sprintf("%s is required", fieldError.Field()),
					})
				case "email":
					logger.Log.Info().
						Str("ip", c.ClientIP()).
						Str("field", fieldError.Field()).
						Msg("Failed registration attempt: invalid email")
					response.Details = append(response.Details, dto.FieldError{
						Field:   fieldError.Field(),
						Code:    "INVALID_EMAIL",
						Message: fmt.Sprintf("Email should be vaild"),
					})
				case "min":
					logger.Log.Info().
						Str("ip", c.ClientIP()).
						Str("field", fieldError.Field()).
						Str("validation_tag", "min").
						Str("required_min", fieldError.Param()).
						Msg("Password length validation failed")
					response.Details = append(response.Details, dto.FieldError{
						Field:   fieldError.Field(),
						Code:    "MIN",
						Message: fmt.Sprintf("Field %s should be longer than %s characters", fieldError.Field(), fieldError.Param()),
					})
				case "max":
					logger.Log.Info().
						Str("ip", c.ClientIP()).
						Str("field", fieldError.Field()).
						Str("validation_tag", "max").
						Str("required_max", fieldError.Param()).
						Msg("Password length validation failed")
					response.Details = append(response.Details, dto.FieldError{
						Field:   fieldError.Field(),
						Code:    "MAX",
						Message: fmt.Sprintf("Field %s should be shorter than %s characters", fieldError.Field(), fieldError.Param()),
					})
				case "type":
					logger.Log.Info().
						Str("ip", c.ClientIP()).
						Str("field", fieldError.Field()).
						Msg("Failed registration attempt: type mismatch")
					response.Details = append(response.Details, dto.FieldError{
						Field:   fieldError.Field(),
						Code:    "TYPE_MISMATCH",
						Message: fmt.Sprintf("%s must be a %s", fieldError.Field(), fieldError.Param()), // e.g., "age must be an integer"
					})
				}
			}
			c.JSON(http.StatusBadRequest, response)
			return
		} else {

			logger.Log.Error().Err(err).Msg("Unhandled exception occured during registration validation.")

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  "INTERNAL_ERROR",
				"text":  "Internal server error",
				"error": err,
			})
			return
		}
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

func (h *AuthHandler) LoginHandler(c *gin.Context) {
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
