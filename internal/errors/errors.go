package myerrors

import "fmt"

type NotFoundError struct {
	Entity string
	ID     string
}

type ValidationError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Entity, e.ID)
}

func (v *ValidationError) Error() string {
	return v.Message
}
