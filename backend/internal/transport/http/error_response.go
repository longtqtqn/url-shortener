package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Authentication & Authorization errors
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeInvalidToken  ErrorCode = "INVALID_TOKEN"
	ErrCodeInvalidAPIKey ErrorCode = "INVALID_API_KEY"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"

	// Validation errors
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField     ErrorCode = "MISSING_FIELD"
	ErrCodeInvalidURL       ErrorCode = "INVALID_URL"
	ErrCodeInvalidEmail     ErrorCode = "INVALID_EMAIL"
	ErrCodePasswordTooShort ErrorCode = "PASSWORD_TOO_SHORT"

	// Resource errors
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict      ErrorCode = "CONFLICT"

	// Business logic errors
	ErrCodeLinkLimitExceeded    ErrorCode = "LINK_LIMIT_EXCEEDED"
	ErrCodeShortCodeUnavailable ErrorCode = "SHORT_CODE_UNAVAILABLE"
	ErrCodeUserAlreadyExists    ErrorCode = "USER_ALREADY_EXISTS"

	// System errors
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Error ValidationErrorDetail `json:"error"`
}

// ValidationErrorDetail contains validation error information
type ValidationErrorDetail struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Fields  []ValidationError `json:"fields,omitempty"`
}

// RespondError sends a standardized error response
func RespondError(ctx *gin.Context, status int, code ErrorCode, message string, details ...string) {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: detail,
		},
	}

	ctx.JSON(status, response)
}

// RespondValidationError sends a validation error response
func RespondValidationError(ctx *gin.Context, message string, fields []ValidationError) {
	response := ValidationErrorResponse{
		Error: ValidationErrorDetail{
			Code:    ErrCodeInvalidInput,
			Message: message,
			Fields:  fields,
		},
	}

	ctx.JSON(http.StatusBadRequest, response)
}

// Common error responses
func RespondUnauthorized(ctx *gin.Context, message ...string) {
	msg := "Authentication required"
	if len(message) > 0 {
		msg = message[0]
	}
	RespondError(ctx, http.StatusUnauthorized, ErrCodeUnauthorized, msg)
}

func RespondForbidden(ctx *gin.Context, message ...string) {
	msg := "Access forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	RespondError(ctx, http.StatusForbidden, ErrCodeForbidden, msg)
}

func RespondNotFound(ctx *gin.Context, resource string) {
	RespondError(ctx, http.StatusNotFound, ErrCodeNotFound, resource+" not found")
}

func RespondConflict(ctx *gin.Context, message string) {
	RespondError(ctx, http.StatusConflict, ErrCodeConflict, message)
}

func RespondInternalError(ctx *gin.Context, message ...string) {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	RespondError(ctx, http.StatusInternalServerError, ErrCodeInternalError, msg)
}

func RespondBadRequest(ctx *gin.Context, message string) {
	RespondError(ctx, http.StatusBadRequest, ErrCodeInvalidInput, message)
}
