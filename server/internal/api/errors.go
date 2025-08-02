package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Error represents a standardized error response format
// @Description Common error response structure for API failures
type Error struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code" example:"UPPERCASE_CODE"`
	Message    string `json:"message" example:"Brief message about the error."`
	Details    any    `json:"details,omitempty"`
}

func (e *Error) SendAndAbort(c *gin.Context) {
	e.addTimingInfo(c)
	c.AbortWithStatusJSON(e.HTTPStatus, e)
}

func (e *Error) SendWithDetailsAndAbort(c *gin.Context, details any) {
	e.Details = details
	e.addTimingInfo(c)
	c.AbortWithStatusJSON(e.HTTPStatus, e)
}

// Helper to add timing information to error
func (e *Error) addTimingInfo(c *gin.Context) {
	if start, exists := c.Get("request_start"); exists {
		if startTime, ok := start.(time.Time); ok {
			c.Header("Request-Latency", time.Since(startTime).String())
		}
	}
}

// General errors
var (
	InternalServerError = &Error{
		HTTPStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
	}

	BadRequest = &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    "Invalid request data",
	}

	InvalidJSON = &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_JSON",
		Message:    "Invalid JSON syntax",
	}

	TypeMismatch = &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "TYPE_MISMATCH",
		Message:    "Field type mismatch",
	}
)

// Auth related
var (
	DuplicateUser = &Error{
		HTTPStatus: http.StatusConflict,
		Code:       "DUPLICATE_USER",
		Message:    "User already exists",
	}

	Unauthorized = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Authentication required",
	}

	InvalidCredentials = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
	}

	Forbidden = &Error{
		HTTPStatus: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    "You don't have access to this data",
	}

	NoToken = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "NO_TOKEN",
		Message:    "Authorization token required",
	}

	BadToken = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "BAD_TOKEN",
		Message:    "Invalid token format",
	}

	InvalidToken = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "INVALID_TOKEN",
		Message:    "Invalid token",
	}

	ExpiredToken = &Error{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "Expired token",
		Message:    "Expired token",
	}

	InvalidQuery = &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_QUERY",
		Message:    "Invalid query parameters",
	}

	NotFound = &Error{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    "Data not found",
	}
)

// Task handler
var (
	InvalidTaskID = &Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "INVALID_ID",
		Message:    "Invalid task ID",
	}
)
