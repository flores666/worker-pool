package service

import (
	"errors"
	"worker-pool/internal/models"
)

type Validator interface {
	Validate(t *models.Task) error
}

type validator struct {
}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(t *models.Task) error {
	if t.Id == "" {
		return errors.New("Empty Id")
	}

	return nil
}
