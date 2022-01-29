package tasks

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]*Task, error)
	Create(task *Task) (*Task, error)
}
type TasksRepo struct {
	tasks *gorm.DB
}

func InitRepo(db *gorm.DB) Repository {
	return &TasksRepo{db.Model(&Task{}).Debug()}
}

func (repo *TasksRepo) GetAll() ([]*Task, error) {
	var tasks []*Task
	err := repo.tasks.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TasksRepo) Create(task *Task) (*Task, error) {
	err := repo.tasks.Create(task).Error
	if err != nil {
		return nil, err
	}

	return task, err
}
