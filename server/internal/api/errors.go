package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	HTTPStatus int    `json:"-"` // Not included in JSON response
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`   // Optional error details
	Timestamp  string `json:"timestamp,omitempty"` // Request start time
	Latency    string `json:"latency,omitempty"`   // Processing duration
}

func (e *APIError) SendAndAbort(c *gin.Context) {
	e.addTimingInfo(c)
	c.AbortWithStatusJSON(e.HTTPStatus, e)
}

func (e *APIError) SendWithDetailsAndAbort(c *gin.Context, details any) {
	e.Details = details
	e.addTimingInfo(c)
	c.AbortWithStatusJSON(e.HTTPStatus, e)
}

// Helper to add timing information to error
func (e *APIError) addTimingInfo(c *gin.Context) {
	if start, exists := c.Get("request_start"); exists {
		if startTime, ok := start.(time.Time); ok {
			e.Timestamp = startTime.Format(time.RFC3339)
			e.Latency = time.Since(startTime).String()
		}
	}
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

	TypeMismatch = &APIError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "TYPE_MISMATCH",
		Message:    "Field type mismatch",
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

	Forbidden = &APIError{
		HTTPStatus: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    "You don't have access to this data",
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

	NotFound = &APIError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    "Data not found",
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
