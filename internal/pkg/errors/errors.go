package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a unique error code
type ErrorCode string

// Error codes for the application
const (
	// Unknown error
	ErrUnknown ErrorCode = "ERR_UNKNOWN"

	// Validation errors (400)
	ErrValidation         ErrorCode = "ERR_VALIDATION"
	ErrMissingField       ErrorCode = "ERR_MISSING_FIELD"
	ErrInvalidField       ErrorCode = "ERR_INVALID_FIELD"
	ErrInvalidCredentials ErrorCode = "ERR_INVALID_CREDENTIALS"

	// Authentication errors (401)
	ErrUnauthenticated ErrorCode = "ERR_UNAUTHENTICATED"
	ErrInvalidToken    ErrorCode = "ERR_INVALID_TOKEN"
	ErrExpiredToken    ErrorCode = "ERR_EXPIRED_TOKEN"
	ErrRevokedToken    ErrorCode = "ERR_REVOKED_TOKEN"

	// Authorization errors (403)
	ErrUnauthorized ErrorCode = "ERR_UNAUTHORIZED"
	ErrForbidden    ErrorCode = "ERR_FORBIDDEN"

	// Not found errors (404)
	ErrNotFound     ErrorCode = "ERR_NOT_FOUND"
	ErrUserNotFound ErrorCode = "ERR_USER_NOT_FOUND"
	ErrTokenNotFound ErrorCode = "ERR_TOKEN_NOT_FOUND"

	// Conflict errors (409)
	ErrConflict      ErrorCode = "ERR_CONFLICT"
	ErrTokenExists   ErrorCode = "ERR_TOKEN_EXISTS"
	ErrUserExists    ErrorCode = "ERR_USER_EXISTS"

	// Internal errors (500)
	ErrInternal       ErrorCode = "ERR_INTERNAL"
	ErrTokenGeneration ErrorCode = "ERR_TOKEN_GENERATION"
	ErrTokenStorage   ErrorCode = "ERR_TOKEN_STORAGE"
	ErrDatabase       ErrorCode = "ERR_DATABASE"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithError sets the underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// Is checks if the error is of the same type
func (e *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return e.Code == appErr.Code
	}
	return false
}

// GetHTTPStatus returns the HTTP status code for the error
func (e *AppError) GetHTTPStatus() int {
	return e.HTTPStatus
}

// GetCode returns the error code
func (e *AppError) GetCode() ErrorCode {
	return e.Code
}

// ToResponse returns a JSON-serializable response
func (e *AppError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error: ErrorDetail{
			Code:    string(e.Code),
			Message: e.Message,
			Details: e.Details,
		},
	}
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Predefined errors with their HTTP status codes
var (
	// Validation errors
	ErrValidationErr = &AppError{
		Code:       ErrValidation,
		Message:    "Validation failed",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrMissingFieldErr = &AppError{
		Code:       ErrMissingField,
		Message:    "Required field is missing",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidFieldErr = &AppError{
		Code:       ErrInvalidField,
		Message:    "Invalid field value",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidCredentialsErr = &AppError{
		Code:       ErrInvalidCredentials,
		Message:    "Invalid username or password",
		HTTPStatus: http.StatusBadRequest,
	}

	// Authentication errors
	ErrUnauthenticatedErr = &AppError{
		Code:       ErrUnauthenticated,
		Message:    "Authentication required",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInvalidTokenErr = &AppError{
		Code:       ErrInvalidToken,
		Message:    "Invalid token",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrExpiredTokenErr = &AppError{
		Code:       ErrExpiredToken,
		Message:    "Token has expired",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrRevokedTokenErr = &AppError{
		Code:       ErrRevokedToken,
		Message:    "Token has been revoked",
		HTTPStatus: http.StatusUnauthorized,
	}

	// Authorization errors
	ErrUnauthorizedErr = &AppError{
		Code:       ErrUnauthorized,
		Message:    "Not authorized to perform this action",
		HTTPStatus: http.StatusForbidden,
	}

	ErrForbiddenErr = &AppError{
		Code:       ErrForbidden,
		Message:    "Access forbidden",
		HTTPStatus: http.StatusForbidden,
	}

	// Not found errors
	ErrNotFoundErr = &AppError{
		Code:       ErrNotFound,
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrUserNotFoundErr = &AppError{
		Code:       ErrUserNotFound,
		Message:    "User not found",
		HTTPStatus: http.StatusNotFound,
	}

	ErrTokenNotFoundErr = &AppError{
		Code:       ErrTokenNotFound,
		Message:    "Token not found",
		HTTPStatus: http.StatusNotFound,
	}

	// Conflict errors
	ErrConflictErr = &AppError{
		Code:       ErrConflict,
		Message:    "Resource conflict",
		HTTPStatus: http.StatusConflict,
	}

	ErrTokenExistsErr = &AppError{
		Code:       ErrTokenExists,
		Message:    "Token already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrUserExistsErr = &AppError{
		Code:       ErrUserExists,
		Message:    "User already exists",
		HTTPStatus: http.StatusConflict,
	}

	// Internal errors
	ErrInternalErr = &AppError{
		Code:       ErrInternal,
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrTokenGenerationErr = &AppError{
		Code:       ErrTokenGeneration,
		Message:    "Failed to generate token",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrTokenStorageErr = &AppError{
		Code:       ErrTokenStorage,
		Message:    "Failed to store token",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrDatabaseErr = &AppError{
		Code:       ErrDatabase,
		Message:    "Database error",
		HTTPStatus: http.StatusInternalServerError,
	}
)

// Common errors for direct use
var (
	ErrUnknownError = errors.New("unknown error")
)

// Helper functions

// NewValidationError creates a validation error with message
func NewValidationError(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &AppError{
		Code:       ErrValidation,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error with resource type
func NewNotFoundError(resource, field string, value string) *AppError {
	return &AppError{
		Code:       ErrNotFound,
		Message:    fmt.Sprintf("%s dengan %s '%s' tidak ditemukan", resource, field, value),
		HTTPStatus: http.StatusNotFound,
	}
}

// NewInternalError creates an internal error with details
func NewInternalError(message string, args ...interface{}) *AppError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &AppError{
		Code:       ErrInternal,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       ErrForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

// NewUnauthenticatedError creates an unauthenticated error
func NewUnauthenticatedError(message string) *AppError {
	return &AppError{
		Code:       ErrUnauthenticated,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}
