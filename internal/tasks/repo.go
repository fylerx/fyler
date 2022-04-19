package tasks

import (
	"fmt"

	"github.com/fylerx/fyler/internal/enum"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]*Task, error)
	GetByID(id uint64) (*Task, error)
	Create(task *Task) (*Task, error)
	SetProgressStatus(task *Task) error
	Failed(task *Task, err error) error
	SetSuccessStatus(task *Task) error
	CloseDBConnection() error
}
type TasksRepo struct {
	db *gorm.DB
}

func InitRepo(db *gorm.DB) Repository {
	return &TasksRepo{db.Session(&gorm.Session{FullSaveAssociations: true}).Preload("Project").Debug()}
}

func (repo *TasksRepo) GetAll() ([]*Task, error) {
	var tasks []*Task
	err := repo.db.Model(&Task{}).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TasksRepo) GetByID(id uint64) (*Task, error) {
	var task *Task
	err := repo.db.Model(&Task{}).Preload("Project.Storage").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return task, err
}

func (repo *TasksRepo) Create(task *Task) (*Task, error) {
	err := repo.db.Model(&Task{}).Create(task).Error
	if err != nil {
		return nil, err
	}

	return task, err
}

func (repo *TasksRepo) Failed(task *Task, err error) error {
	return repo.db.Model(task).
		Select("status", "error").
		Updates(Task{Status: enum.StatusFailed, Error: err.Error()}).
		Error
}

func (repo *TasksRepo) SetSuccessStatus(task *Task) error {
	task.Status = enum.StatusSuccess
	return repo.db.Select("Status", "Conversion").Updates(task).Error
}

func (repo *TasksRepo) SetProgressStatus(task *Task) error {
	return repo.db.Model(task).Select("status").Update("status", enum.StatusProgress).Error
}

func (repo *TasksRepo) CloseDBConnection() error {
	dbConn, err := repo.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get psql connection: %w", err)
	}

	return dbConn.Close()
}
