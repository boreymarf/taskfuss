package security

import (
	"os"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashCost, err := strconv.Atoi(os.Getenv("HASH_COST"))

	if err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to get hash cost, resetting to 10")
		hashCost = 10
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return "", err
	}

	return string(passwordHash), nil
}
