package entity

import "strings"

type TypeError string

const (
	ENTITY_ERROR      TypeError = "error on validate entity"
	NOT_FOUND_ERROR   TypeError = "not found entity"
	INTERNAL_ERROR    TypeError = "internal error"
	NOT_ALLOWED_ERROR TypeError = "not allowed action"
)

type ErrorHandler struct {
	Messages  []string
	TypeError TypeError
}

func (e *ErrorHandler) Add(message string) {
	e.Messages = append(e.Messages, message)
}

func (e *ErrorHandler) Error() string {
	return strings.Join(e.Messages, ", ")
}

func (e *ErrorHandler) GetTypeError() TypeError {
	return e.TypeError
}
