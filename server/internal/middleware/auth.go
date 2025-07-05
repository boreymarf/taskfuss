package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/gin-gonic/gin"
)

func Auth(userRepo *db.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn().Msg("Auth attempt with no token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
				Code:    "NO_TOKEN",
				Message: "Authorization token required",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Log.Warn().Msg("Auth attempt with a bad token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
				Code:    "BAD_TOKEN",
				Message: "Invalid token format",
			})
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			logger.Log.Error().Msg("Cannot generate JWT token since the secret is not set in the .env!")
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.InternalError{})
			return
		}

		claims, err := security.VerifyToken(tokenParts[1], []byte(secret))

		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidToken) {
				logger.Log.Warn().Msg("Auth attempt with a invalid token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
					Code:    "INVALID_TOKEN",
					Message: "Invalid token",
				})
			} else if errors.Is(err, apperrors.ErrUnexpectedSigningMethod) {
				logger.Log.Warn().Msg("Auth attempt with a bad token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
					Code:    "BAD_TOKEN",
					Message: "Invalid token format",
				})
			} else if errors.Is(err, apperrors.ErrTokenExpired) {
				logger.Log.Warn().Msg("Auth attempt with a expired token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
					Code:    "EXPIRED_TOKEN",
					Message: "Expired token",
				})
			} else {
				logger.Log.Error().Err(err).Msg("Auth attempt failed")
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.InternalError{
					Code:    "INTERNAL_ERROR",
					Message: "Internal error",
				})
			}
			return
		}

		if claims.UserID <= 0 {
			logger.Log.Warn().Int64("user_id", claims.UserID).Msg("Auth attempt failed, invalid user id in the claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
				Code:    "INVALID_USERID",
				Message: "User ID must be a positive integer",
			})
			return
		}

		exists, err := userRepo.Exists(claims.UserID)
		if err != nil {
			logger.Log.Error().Err(err).Send()
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.InternalError{
				Code:    "INTERNAL_ERROR",
				Message: "Internal error",
			})
		}
		if !exists {
			logger.Log.Warn().Int64("user_id", claims.UserID).Msg("Auth attempt failed, user does not exist")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.GenericError{
				Code:    "INVALID_TOKEN",
				Message: "Invalid token",
			})
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
