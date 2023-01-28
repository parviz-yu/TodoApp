package store

import "github.com/pyuldashev912/todoapp/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}

type TaskRepository interface {
	Create(*model.Task) error
	Delete(int, int) error
	Done(int, int) error
	GetAll(int) ([]*model.Task, error)
	GetBool(int, bool) ([]*model.Task, error)
	GetById(int, int) (*model.Task, error)
}
