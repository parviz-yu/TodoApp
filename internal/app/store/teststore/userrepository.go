package teststore

import (
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type UserRepository struct {
	users map[int]*model.User
}

func (r *UserRepository) Create(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.EncryptPassword(); err != nil {
		return err
	}

	user.ID = len(r.users) + 1
	r.users[user.ID] = user

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, store.ErrNoRecordsInTable
}
