package sqlstore_test

import (
	"testing"
	"time"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
	"github.com/pyuldashev912/todoapp/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("tasks")

	s := sqlstore.New(db)
	task := model.TestTask(t)
	err := s.Task().Create(task)
	assert.NoError(t, err)
}

func TestTaskRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("tasks")

	s := sqlstore.New(db)
	_, err := s.Task().GetAll(1)
	assert.EqualError(t, err, store.ErrNoRecordsInTable.Error())

	task := model.TestTask(t)
	s.Task().Create(task)
	s.Task().Create(task)
	result, err := s.Task().GetAll(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestTaskRepository_GetBool(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("tasks")

	s := sqlstore.New(db)
	_, err := s.Task().GetBool(1, false)
	assert.EqualError(t, err, store.ErrNoRecordsInTable.Error())

	task := model.TestTask(t)
	s.Task().Create(task)
	s.Task().Create(task)
	s.Task().Create(&model.Task{
		UserID: 2, Title: "Check", Description: "Some text", CreationDate: time.Now(),
	})
	result, err := s.Task().GetBool(2, false)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
}
