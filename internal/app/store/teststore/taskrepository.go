package teststore

import (
	"github.com/pyuldashev912/todoapp/internal/app/model"
	"github.com/pyuldashev912/todoapp/internal/app/store"
)

type TaskRepository struct {
	tasks map[int]*model.Task
}

func (r *TaskRepository) Create(task *model.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	task.ID = len(r.tasks) + 1
	r.tasks[task.ID] = task

	return nil
}

func (r *TaskRepository) Delete(userId int, taskId int) error {
	targetTaskId, err := getKeyFromMap(r.tasks, userId, taskId)
	if err != nil {
		return err
	}

	delete(r.tasks, targetTaskId)
	return nil
}

func (r *TaskRepository) Done(userId int, taskId int) error {
	targetTaskId, err := getKeyFromMap(r.tasks, userId, taskId)
	if err != nil {
		return err
	}

	r.tasks[targetTaskId].Done = true
	return nil
}

func (r *TaskRepository) GetById(userId int, taskId int) (*model.Task, error) {
	tagrgetid, err := getKeyFromMap(r.tasks, userId, taskId)
	if err != nil {
		return nil, err
	}

	return r.tasks[tagrgetid], nil
}

func (r *TaskRepository) GetBool(userId int, done bool) ([]*model.Task, error) {
	var tasks []*model.Task
	for _, task := range r.tasks {
		if task.UserID == userId && task.Done == done {
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		return nil, store.ErrNoRecordsInTable
	}

	return tasks, nil
}

func (r *TaskRepository) GetAll(userId int) ([]*model.Task, error) {
	var tasks []*model.Task
	for _, task := range r.tasks {
		if task.UserID == userId {
			tasks = append(tasks, task)
		}
	}

	if len(tasks) == 0 {
		return nil, store.ErrNoRecordsInTable
	}

	return tasks, nil
}

func getKeyFromMap(tasks map[int]*model.Task, userId int, taskId int) (int, error) {
	var targetTaskId int
	for k, task := range tasks {
		if task.UserID == userId && task.ID == taskId {
			targetTaskId = k
			break
		}
	}

	if _, ok := tasks[targetTaskId]; !ok {
		return 0, store.ErrInvalidTaskId
	}

	return targetTaskId, nil
}
