package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

const (
	ErrCodeDatabase        = "DATABASE_ERROR"
	ErrCodeKafka           = "KAFKA_ERROR"
	ErrCodeFCM             = "FCM_ERROR"
	ErrCodeInvalidPayload  = "INVALID_PAYLOAD"
	ErrCodeConfiguration   = "CONFIGURATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInternal        = "INTERNAL_ERROR"
)

func NewAppError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(ErrCodeDatabase, message, err)
}

func NewKafkaError(message string, err error) *AppError {
	return NewAppError(ErrCodeKafka, message, err)
}

func NewFCMError(message string, err error) *AppError {
	return NewAppError(ErrCodeFCM, message, err)
}

func NewInvalidPayloadError(message string, err error) *AppError {
	return NewAppError(ErrCodeInvalidPayload, message, err)
}
