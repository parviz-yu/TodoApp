package model

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) *User {
	return &User{
		Name:     "Parviz",
		Email:    "user@user.com",
		Password: "Password",
	}
}

func TestTask(t *testing.T) *Task {
	return &Task{
		UserID:       1,
		Title:        "Create TodoApp",
		Description:  "Add new features",
		Done:         false,
		CreationDate: time.Now().String(),
	}
}
