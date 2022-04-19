package tasks

import (
	"time"

	"github.com/fylerx/fyler/internal/conversions"
	"github.com/fylerx/fyler/internal/enum"
	"github.com/fylerx/fyler/internal/projects"
)

type Task struct {
	ID         uint64                  `gorm:"primarykey" json:"id"`
	ProjectID  uint64                  `json:"project_id"`
	Project    *projects.Project       `json:"-"`
	Status     enum.Status             `gorm:"type:status;default:'queued'" json:"status"`
	TaskType   enum.TaskType           `gorm:"type:task_type" json:"task_type"`
	FilePath   string                  `json:"file_path"`
	Error      string                  `json:"error"`
	Conversion *conversions.Conversion `json:"conversion"`
	CreatedAt  time.Time               `json:"created_at"`
	UpdatedAt  time.Time               `json:"updated_at"`
}
