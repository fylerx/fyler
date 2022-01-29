package tasks

import (
	"time"

	"github.com/fylerx/fyler/internal/enum"
	"github.com/fylerx/fyler/internal/projects"
)

type Task struct {
	ID        uint64 `gorm:"primarykey"`
	ProjectID uint32
	Project   projects.Project
	Status    enum.Status   `gorm:"type:status;default:'queued'"`
	TaskType  enum.TaskType `gorm:"type:task_type"`
	URL       string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
