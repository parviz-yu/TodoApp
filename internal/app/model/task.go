package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Task struct {
	ID           int    `json:"id"`
	UserID       int    `json:"-"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Done         bool   `json:"done"`
	CreationDate string `json:"creation_date"`
}

func (t *Task) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Title, validation.Required),
		validation.Field(&t.Description, validation.Required),
	)
}
