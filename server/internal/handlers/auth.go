package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/api"
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
	userRepo *db.Users
}

func InitAuthHandler(userRepo *db.Users) (*AuthHandler, error) {
	return &AuthHandler{userRepo: userRepo}, nil
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account and return a JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param RegisterRequest body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.RegisterResponse "Successfully registered"
// @Failure 400 {object} api.Error "Invalid request format (code: BAD_REQUEST) or username/email already exists (code: DUPLICATE_USER)"
// @Failure 500 {object} api.Error "Internal server error (code: INTERNAL_ERROR)
// @Router /auth/register [post]
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

		api.InternalServerError.SendAndAbort(c)
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
			api.DuplicateUser.SendAndAbort(c)
			return
		}

		logger.Log.Error().
			Err(err).
			Str("name", req.Username).
			Str("email", req.Email).
			Msg("Failed to create new user in the database")

		api.InternalServerError.SendAndAbort(c)
		return
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := security.CreateToken(user.ID, []byte(secret), time.Hour*24)
	if err != nil {
		api.InternalServerError.SendAndAbort(c)
		return
	}

	api.Created(c, dto.RegisterResponse{
		User: dto.User{
			Id:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		AuthToken: token,
	})

}

// PingHandler godoc
// @Summary Server health check
// @Description Returns "pong" if the server is running
// @Tags service
// @Produce json
// @Success 200 {object} dto.PongResponse "Server is running"
// @Router /api/ping [get]

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate user credentials and return a JWT token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param   LoginRequest  body  dto.LoginRequest  true  "Login credentials"
// @Success 200 {object}  dto.LoginResponse  "Successfully authenticated"
// @Failure 400 {object}  api.Error                          "Invalid request format"
// @Failure 401 {object}  api.Error                          "Invalid credentials"
// @Failure 500 {object}  api.Error                          "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {

			logger.Log.Error().Err(syntaxErr).Send()

			api.InvalidJSON.SendAndAbort(c)
			return
		}

		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {

			logger.Log.Error().
				Str("field", typeErr.Field).
				Str("expected_type", typeErr.Type.String()).
				Msg("Type mismatch")

			details := api.FieldErrorDetail{
				Field:    typeErr.Field,
				Expected: typeErr.Type.String(),
				Message:  fmt.Sprintf("Field '%s' must be a %s", typeErr.Field, typeErr.Type),
			}
			api.TypeMismatch.SendWithDetailsAndAbort(c, details)
			return
		}
	}

	var user models.User

	err := h.userRepo.GetUserByEmail(req.Email, &user)

	logger.Log.Debug().Str("email", user.Email).Str("hash", user.PasswordHash).Send()

	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: user does not exists")
			api.InvalidCredentials.SendAndAbort(c)
			return
		} else {
			logger.Log.Error().Str("email", req.Email).Err(err).Msg("Failed login attempt: internal server error")
			api.InternalServerError.SendAndAbort(c)
			return
		}
	}

	if req.Password == "" {

		logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: empty password")
		api.InvalidCredentials.SendAndAbort(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {

		logger.Log.Warn().Str("email", req.Email).Err(err).Msg("Failed login attempt: incorrect password")
		api.InvalidCredentials.SendAndAbort(c)
		return
	}

	secret := os.Getenv("JWT_SECRET")

	token, err := security.CreateToken(user.ID, []byte(secret), time.Hour*24)
	if err != nil {
		api.InternalServerError.SendAndAbort(c)
		return
	}

	api.Success(c, dto.RegisterResponse{
		User: dto.User{
			Id:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
		},
		AuthToken: token,
	})
}
