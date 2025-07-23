// internal/pkg/errors.go
package pkg

import "fmt"

// Custom error types

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(format string, a ...interface{}) error {
	return &NotFoundError{Message: fmt.Sprintf(format, a...)}
}

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

func NewInvalidInputError(format string, a ...interface{}) error {
	return &InvalidInputError{Message: fmt.Sprintf(format, a...)}
}

type InternalServerError struct {
	Message string
	Err     error // Original error
}

func (e *InternalServerError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewInternalServerError(err error, format string, a ...interface{}) error {
	return &InternalServerError{Message: fmt.Sprintf(format, a...), Err: err}
}
