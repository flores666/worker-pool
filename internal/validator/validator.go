package validator

import (
	"errors"
	"worker-pool/internal/tasks"
)

type Validator interface {
	Validate(t tasks.Task) error
}

type validator struct {
}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(t tasks.Task) error {
	if t.Id == "" {
		return errors.New("Empty Id")
	}

	if t.Value == "" {
		return errors.New("Empty Value")
	}

	return nil
}
