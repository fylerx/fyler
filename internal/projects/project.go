package projects

import (
	"context"
	"fmt"
	"time"

	"github.com/fylerx/fyler/internal/storages"
)

type Project struct {
	ID        uint64           `gorm:"primarykey" json:"id"`
	Name      string           `json:"name"`
	APIKey    string           `json:"apikey"`
	Storage   storages.Storage `json:"-"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
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
