package apperrors

import "fmt"

type ValidationError struct {
	Code    string
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("error %s: %s", e.Code, e.Message)
}

func NewValidationError(code string, field string, message string) *ValidationError {
	return &ValidationError{
		Code:    code,
		Field:   field,
		Message: message,
	}
}
