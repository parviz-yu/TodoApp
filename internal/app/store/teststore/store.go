package teststore

import (
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
	taskRepository *TaskRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		users: make(map[int]*model.User),
	}

	return s.userRepository
}

func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		tasks: make(map[int]*model.Task),
	}

	return s.taskRepository
}
