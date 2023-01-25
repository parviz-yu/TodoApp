package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Task struct {
	ID           int
	UserID       int
	Title        string
	Description  string
	Done         bool
	CreationDate time.Time
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Title, validation.Required),
		validation.Field(&t.Description, validation.Required),
	)
}
