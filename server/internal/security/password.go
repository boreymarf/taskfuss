package security

import (
	"fmt"
	"os"
	"strconv"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hashCost, err := strconv.Atoi(os.Getenv("HASH_COST"))
	if err != nil || hashCost < 4 || hashCost > 31 {
		logger.Log.Warn().Msgf("Invalid HASH_COST (%d). Using default value 10", hashCost)
		hashCost = 10
	}

	if password == "" {
		logger.Log.Error().Msg("Password cannot be empty")
		return "", fmt.Errorf("password is empty")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		return "", err
	}

	return string(passwordHash), nil
}
