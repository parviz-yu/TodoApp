package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	taskRepository *TaskRepository
}

// NewStore returns a new instance of store.
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User returns a userRepository. It is used to interact with the repository from the outside.
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

// Task returns a taskRepository. It is used to interact with the repository from the outside.
func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
	}

	return s.taskRepository
}
