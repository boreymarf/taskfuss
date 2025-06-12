package security

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserID    int64  `json:"user_id"`
	Usernamse string `json:"username"`
	jwt.RegisteredClaims
}
