package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *APIError) Send(c *gin.Context) {
	c.JSON(e.HTTPStatus, gin.H{
		"code":    e.Code,
		"message": e.Message,
	})
}

func (e *APIError) SendWithDetails(c *gin.Context, details any) {
	c.JSON(e.HTTPStatus, gin.H{
		"code":    e.Code,
		"text":    e.Message,
		"details": details,
	})
}

// General errors
var (
	InternalServerError = &APIError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
	}

	BadRequest = &APIError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    "Invalid request data",
	}

	InvalidJSON = &APIError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_JSON",
		Message:    "Invalid JSON syntax",
	}
)

// Auth related
var (
	DuplicateUser = &APIError{
		HTTPStatus: http.StatusConflict,
		Code:       "DUPLICATE_USER",
		Message:    "User already exists",
	}

	Unauthorized = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
	}

	InvalidCredentials = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
	}
)
