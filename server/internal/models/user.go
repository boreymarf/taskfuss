package models

import "time"

type UserRole int

const (
	RoleAdmin UserRole = iota
	RoleUser
	RoleGuest
)

type UserContext struct {
	ID   int64    `db:"id"`
	Role UserRole `db:"role"`
}

type User struct {
	ID           int64     `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Role         UserRole  `db:"role"`
}
