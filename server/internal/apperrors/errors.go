package apperrors

import "errors"

var (
	ErrDuplicate = errors.New("duplicate_entry")
	ErrNotFound  = errors.New("not_found")
)
