package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=40"`
}

type RegisterResponse struct {
	UserID    int64  `json:"userId"`
	Status    string `json:"status"`
	AuthToken string `json:"auth_token"`
	ExpiresAt int64  `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	UserID    int64  `json:"userId"`
	Status    string `json:"status"`
	AuthToken string `json:"auth_token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ValidationError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details"`
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
