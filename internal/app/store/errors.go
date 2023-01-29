package store

import "errors"

var (
	ErrNoRecordsInTable = errors.New("no records in table")
	ErrInvalidTaskId    = errors.New("invalid task id")
)
