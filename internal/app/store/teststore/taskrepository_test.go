package teststore_test

import (
	"testing"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
	"github.com/pyuldashev912/todoapp/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	s := teststore.New()
	task := model.TestTask(t)
	err := s.Task().Create(task)
	assert.NoError(t, err)
}

func TestTaskRepository_Delete(t *testing.T) {
	s := teststore.New()
	err := s.Task().Delete(1, 1)
	assert.EqualError(t, err, store.ErrInvalidTaskId.Error())

	task := model.TestTask(t)
	s.Task().Create(task)
	err = s.Task().Delete(task.UserID, 1)
	assert.NoError(t, err)
}

func TestTaskRepository_GetById(t *testing.T) {
	s := teststore.New()
	task := model.TestTask(t)
	s.Task().Create(task)
	_, err := s.Task().GetById(task.UserID, 1)
	assert.NoError(t, err)
}

func TestTaskRepository_Done(t *testing.T) {
	s := teststore.New()
	task := model.TestTask(t)
	s.Task().Create(task)
	err := s.Task().Done(task.UserID, 1)
	res, _ := s.Task().GetById(task.UserID, 1)
	assert.NoError(t, err)
	assert.Equal(t, task.Done, res.Done)
}

func TestTaskRepository_GetBool(t *testing.T) {
	s := teststore.New()
	task := model.TestTask(t)
	s.Task().Create(task)
	_, err := s.Task().GetBool(task.UserID, true)
	assert.EqualError(t, err, store.ErrNoRecordsInTable.Error())

	res, err := s.Task().GetBool(task.UserID, false)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestTaskRepository_GetAll(t *testing.T) {
	s := teststore.New()
	_, err := s.Task().GetAll(5)
	assert.EqualError(t, err, store.ErrNoRecordsInTable.Error())

	task := model.TestTask(t)
	s.Task().Create(task)
	_, err = s.Task().GetAll(task.UserID)
	assert.NoError(t, err)
}
