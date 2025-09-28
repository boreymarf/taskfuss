package models

import "time"

type UserRole int

const (
	RoleAdmin UserRole = iota
	RoleUser
	RoleGuest
)

type UserContext struct {
	ID   int64
	Role UserRole
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash,omitempty"` // TODO: Hash password later
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
