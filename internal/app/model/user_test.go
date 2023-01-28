package model_test

import (
	"testing"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""

				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "invalid"

				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""

				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "123"

				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestUser_EncryptPassword(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.EncryptPassword())
	assert.NotEmpty(t, u.EncryptedPassword)
}

func TestUser_ComparePassword(t *testing.T) {
	u := model.TestUser(t)
	u.EncryptPassword()
	result := u.ComparePassword(u.Password)
	assert.True(t, result)
}
func TestUser_Sanitize(t *testing.T) {
	u := model.TestUser(t)
	u.Sanitize()
	assert.Empty(t, u.Password)
}
