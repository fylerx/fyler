package tasks

import (
	"time"

	"github.com/fylerx/fyler/internal/enum"
	"github.com/fylerx/fyler/internal/projects"
)

type Task struct {
	ID        uint64           `gorm:"primarykey" json:"id"`
	ProjectID uint32           `json:"project_id"`
	Project   projects.Project `json:"-"`
	Status    enum.Status      `gorm:"type:status;default:'queued'" json:"status"`
	TaskType  enum.TaskType    `gorm:"type:task_type" json:"task_type"`
	URL       string           `json:"url"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
