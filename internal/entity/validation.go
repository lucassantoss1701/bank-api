package entity

import "strings"

type ValidationError struct {
	Messages []string
}

func (v *ValidationError) Add(message string) {
	v.Messages = append(v.Messages, message)
}

func (v *ValidationError) Error() string {
	return strings.Join(v.Messages, ", ")
}
