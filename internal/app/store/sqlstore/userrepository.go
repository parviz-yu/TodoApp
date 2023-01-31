package sqlstore

import (
	"database/sql"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type UserRepository struct {
	store *Store
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	// Validate users field
	if err := user.Validate(); err != nil {
		return err
	}

	// Encrypt password before creating a user
	if err := user.EncryptPassword(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		`INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Email, user.EncryptedPassword,
	).Scan(&user.ID)
}

// FindByEmail returns a user with appropriate email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id, name, email, encrypted_password FROM users WHERE email=$1`, email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword); err != nil {
		return nil, err
	}

	return user, nil
}

// FindById returns a user with appropriate id
func (r *UserRepository) FindById(userId int) (*model.User, error) {
	user := &model.User{}
	if err := r.store.db.QueryRow(
		`SELECT id, name, email FROM users WHERE id=$1`, userId,
	).Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNoRecordsInTable
		}

		return nil, err
	}

	return user, nil
}
