package sqlstore

import (
	"database/sql"

	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type TaskRepository struct {
	store *Store
}

// Create creates a new task
func (r *TaskRepository) Create(task *model.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(`
	INSERT INTO tasks (user_id, title, description, done, creation_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		task.UserID, task.Title, task.Description, task.Done, task.CreationDate,
	).Scan(&task.ID)
}

// GetAll gets all User's tasks
func (r *TaskRepository) GetAll(userId int) ([]*model.Task, error) {
	return r.getUnderHood(userId)
}

// GetBool gets all User's tasks that are completed or not completed
func (r *TaskRepository) GetBool(userId int, status bool) ([]*model.Task, error) {
	return r.getUnderHood(userId, status)
}

// Slice of empty interface allows to make more complex database queries
func (r *TaskRepository) getUnderHood(values ...interface{}) ([]*model.Task, error) {
	var rows *sql.Rows
	var err error

	// When we need to get all concrete User's tasks
	if len(values) == 1 {
		rows, err = r.store.db.Query(
			"SELECT * FROM tasks WHERE user_id=$1", values[0].(int),
		)
	}

	// When we need to get all completed/not completed concrete User's tasks
	if len(values) == 2 {
		rows, err = r.store.db.Query(
			"SELECT * FROM tasks WHERE user_id=$1 and done=$2",
			values[0].(int), values[1].(bool),
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*model.Task, 0, 5)
	for rows.Next() {
		s := &model.Task{}
		err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.Description, &s.Done, &s.CreationDate)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, store.ErrNoRecordsInTable
	}

	return tasks, nil
}
