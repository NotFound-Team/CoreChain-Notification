package errors

import "fmt"

// AppError represents application-specific errors
type AppError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeDatabase        = "DATABASE_ERROR"
	ErrCodeKafka           = "KAFKA_ERROR"
	ErrCodeFCM             = "FCM_ERROR"
	ErrCodeInvalidPayload  = "INVALID_PAYLOAD"
	ErrCodeConfiguration   = "CONFIGURATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInternal        = "INTERNAL_ERROR"
)

// NewAppError creates a new application error
func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewDatabaseError creates a database-related error
func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(ErrCodeDatabase, message, err)
}

// NewKafkaError creates a Kafka-related error
func NewKafkaError(message string, err error) *AppError {
	return NewAppError(ErrCodeKafka, message, err)
}

// NewFCMError creates an FCM-related error
func NewFCMError(message string, err error) *AppError {
	return NewAppError(ErrCodeFCM, message, err)
}

// NewInvalidPayloadError creates an invalid payload error
func NewInvalidPayloadError(message string, err error) *AppError {
	return NewAppError(ErrCodeInvalidPayload, message, err)
}
