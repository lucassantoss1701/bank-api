package entity

import "strings"

type TypeError string

const (
	ENTITY_ERROR       TypeError = "error on validate entity"
	NOT_FOUND_ERROR    TypeError = "not found"
	INTERNAL_ERROR     TypeError = "internal error"
	CONFLICT_ERROR     TypeError = "conflict error"
	NOT_ALLOWED_ERROR  TypeError = "not allowed"
	UNAUTHORIZED_ERROR TypeError = "unauthorized"
	BAD_REQUEST        TypeError = "bad request"
)

type ErrorHandler struct {
	Messages  []string
	TypeError TypeError
}

func NewErrorHandler(typeError TypeError) *ErrorHandler {
	return &ErrorHandler{
		TypeError: typeError,
	}
}

func (e *ErrorHandler) Add(message string) *ErrorHandler {
	e.Messages = append(e.Messages, message)
	return e
}

func (e *ErrorHandler) Error() string {
	return strings.Join(e.Messages, ", ")
}

func (e *ErrorHandler) GetTypeError() TypeError {
	return e.TypeError
}
