package apperrors

import "errors"

var (
	ErrInvalidToken            = errors.New("invalid_token")
	ErrUnexpectedSigningMethod = errors.New("unexpected_signing_method")
	ErrTokenExpired            = errors.New("token_expired")
)
