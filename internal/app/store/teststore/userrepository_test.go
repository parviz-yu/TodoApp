package teststore_test

import (
	"testing"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
	"github.com/pyuldashev912/todoapp/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	store := teststore.New()
	user := model.TestUser(t)
	err := store.User().Create(user)
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	_, err := s.User().FindByEmail("user@password.com")
	assert.EqualError(t, err, store.ErrNoRecordsInTable.Error())

	user := model.TestUser(t)
	s.User().Create(user)
	res, err := s.User().FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	s.User().Create(u1)

	u2, err := s.User().FindById(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)

}
