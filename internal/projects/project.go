package projects

import (
	"context"
	"fmt"
	"time"
)

type Project struct {
	ID        uint32    `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"apikey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type projectKey string

var ProjectKey projectKey = "currentProject"

func ProjectFromContext(ctx context.Context) (*Project, error) {
	project, ok := ctx.Value(ProjectKey).(*Project)
	if !ok || project == nil {
		return nil, fmt.Errorf("invalid project")
	}
	return project, nil
}
