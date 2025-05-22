package security

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateToken(userID int64, secret []byte, expiresIn time.Duration) (string, error) {

	expiresAt := time.Now().Add(expiresIn)

	claims := CustomClaims{
		UserID: strconv.Itoa(int(userID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "taskfuss",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	logger.Log.Debug().Str("token", token.Raw).Int("expires in", int(expiresIn)).Msg("Created token")

	return tokenString, nil
}

func VerifyToken(tokenString string, secret []byte) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (any, error) {
			// Validate the signing method (e.g., HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil // Return the secret key used for signing
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		logger.Log.Debug().Str("userID", claims.UserID).Msg("Token verified")

		return claims, nil
	}

	return nil, errors.New("invalid token")
}
