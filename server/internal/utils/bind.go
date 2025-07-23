package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/boreymarf/task-fuss/server/internal/api"
	"github.com/boreymarf/task-fuss/server/internal/dto"
	"github.com/boreymarf/task-fuss/server/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func HandleBindingError(c *gin.Context, err error) {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		logger.Log.Warn().Err(syntaxErr).Msg("Incorrect type in the JSON")
		api.InvalidJSON.SendAndAbort(c)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		logger.Log.Warn().Err(err).Msg("Unexpected EOF in the JSON request")
		api.InvalidJSON.SendAndAbort(c)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {

		logger.Log.Error().
			Str("field", typeErr.Field).
			Str("expected_type", typeErr.Type.String()).
			Msg("Type mismatch")

		c.JSON(http.StatusBadRequest, gin.H{
			"code": "TYPE_MISMATCH",
			"text": fmt.Sprintf("Field '%s' must be a %s", typeErr.Field, typeErr.Type),
		})
		return
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {

		logger.Log.Warn().Err(validationErrors).Send()

		var response dto.ValidationError = dto.ValidationError{
			Code:    "VALIDATION_FAILED",
			Message: "Validation failed",
			Details: make([]dto.FieldError, 0),
		}

		for _, fieldError := range validationErrors {

			switch fieldError.Tag() {
			case "required":
				logger.Log.Info().
					Str("ip", c.ClientIP()).
					Str("field", fieldError.Field()).
					Msg("Failed registration attempt: missing required field")
				response.Details = append(response.Details, dto.FieldError{
					Field:   fieldError.Field(),
					Code:    "REQUIRED",
					Message: fmt.Sprintf("%s is required", fieldError.Field()),
				})
			case "email":
				logger.Log.Info().
					Str("ip", c.ClientIP()).
					Str("field", fieldError.Field()).
					Msg("Failed registration attempt: invalid email")
				response.Details = append(response.Details, dto.FieldError{
					Field:   fieldError.Field(),
					Code:    "INVALID_EMAIL",
					Message: "Email should be vaild",
				})
			case "min":
				logger.Log.Info().
					Str("ip", c.ClientIP()).
					Str("field", fieldError.Field()).
					Str("validation_tag", "min").
					Str("required_min", fieldError.Param()).
					Msg("Password length validation failed")
				response.Details = append(response.Details, dto.FieldError{
					Field:   fieldError.Field(),
					Code:    "MIN",
					Message: fmt.Sprintf("Field %s should be longer than %s characters", fieldError.Field(), fieldError.Param()),
				})
			case "max":
				logger.Log.Info().
					Str("ip", c.ClientIP()).
					Str("field", fieldError.Field()).
					Str("validation_tag", "max").
					Str("required_max", fieldError.Param()).
					Msg("Password length validation failed")
				response.Details = append(response.Details, dto.FieldError{
					Field:   fieldError.Field(),
					Code:    "MAX",
					Message: fmt.Sprintf("Field %s should be shorter than %s characters", fieldError.Field(), fieldError.Param()),
				})
			case "type":
				logger.Log.Info().
					Str("ip", c.ClientIP()).
					Str("field", fieldError.Field()).
					Msg("Failed registration attempt: type mismatch")
				response.Details = append(response.Details, dto.FieldError{
					Field:   fieldError.Field(),
					Code:    "TYPE_MISMATCH",
					Message: fmt.Sprintf("%s must be a %s", fieldError.Field(), fieldError.Param()), // e.g., "age must be an integer"
				})
			}
		}
		c.JSON(http.StatusBadRequest, response)
		return
	} else {

		logger.Log.Error().Err(err).Msg("Unhandled exception occured during registration validation.")

		api.InternalServerError.SendAndAbort(c)
		return
	}
}
