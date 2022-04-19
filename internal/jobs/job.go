package jobs

import "github.com/fylerx/fyler/internal/enum"

type Job struct {
	TaskID   uint64        `json:"task_id"`
	TaskType enum.TaskType `json:"task_type"`
}
