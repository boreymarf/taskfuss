package middleware

import (
	"errors"
	"os"
	"strings"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/apperrors"
	"github.com/boreymarf/task-fuss/server/internal/db"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/boreymarf/task-fuss/server/internal/security"
	"github.com/gin-gonic/gin"
)

func Auth(userRepo *db.Users) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Log.Warn().Msg("Auth attempt with no token")
			api.NoToken.SendAndAbort(c)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Log.Warn().Msg("Auth attempt with a bad token")
			api.BadToken.SendAndAbort(c)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			logger.Log.Error().Msg("Cannot generate JWT token since the secret is not set in the .env!")
			api.InternalServerError.SendAndAbort(c)
			return
		}

		claims, err := security.VerifyToken(tokenParts[1], []byte(secret))

		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidToken) {
				logger.Log.Warn().Msg("Auth attempt with a invalid token")
				api.InvalidToken.SendAndAbort(c)
			} else if errors.Is(err, apperrors.ErrUnexpectedSigningMethod) {
				logger.Log.Warn().Msg("Auth attempt with a bad token")
				api.BadToken.SendAndAbort(c)
			} else if errors.Is(err, apperrors.ErrTokenExpired) {
				logger.Log.Warn().Msg("Auth attempt with a expired token")
				api.ExpiredToken.SendAndAbort(c)
			} else {
				logger.Log.Error().Err(err).Msg("Auth attempt failed")
				api.InternalServerError.SendAndAbort(c)
			}
			return
		}

		if claims.UserID <= 0 {
			logger.Log.Warn().Int64("user_id", claims.UserID).Msg("Auth attempt failed, invalid user id in the claims")
			api.InternalServerError.SendAndAbort(c)
			return
		}

		exists, err := userRepo.Exists(claims.UserID)
		if err != nil {
			logger.Log.Error().Err(err).Send()
			api.InternalServerError.SendAndAbort(c)
		}
		if !exists {
			logger.Log.Warn().Int64("user_id", claims.UserID).Msg("Auth attempt failed, user does not exist")
			api.InvalidToken.SendAndAbort(c)
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
