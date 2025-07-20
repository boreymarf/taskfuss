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

func (e *APIError) SendAndAbort(c *gin.Context) {
	c.AbortWithStatusJSON(e.HTTPStatus, gin.H{
		"code":    e.Code,
		"message": e.Message,
	})
}

func (e *APIError) SendWithDetailsAndAbort(c *gin.Context, details any) {
	c.AbortWithStatusJSON(e.HTTPStatus, gin.H{
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

	NoToken = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "NO_TOKEN",
		Message:    "Authorization token required",
	}

	BadToken = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "BAD_TOKEN",
		Message:    "Invalid token format",
	}

	InvalidToken = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "INVALID_TOKEN",
		Message:    "Invalid token",
	}

	ExpiredToken = &APIError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "Expired token",
		Message:    "Expired token",
	}

	InvalidQuery = &APIError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_QUERY",
		Message:    "Invalid query parameters",
	}
)

// Task handler
var (
	InvalidTaskID = &APIError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_ID",
		Message:    "Invalid task ID",
	}
)
