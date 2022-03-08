package projects_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/fylerx/fyler/internal/projects"
)

func TestProjectFromContext(t *testing.T) {
	prj := &projects.Project{ID: 555, Name: "New Project", APIKey: "api-key"}
	tests := []struct {
		name    string
		method  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "valid project",
			method:  prj,
			want:    prj,
			wantErr: false,
		},
		{
			name:    "project is nil",
			method:  nil,
			want:    prj,
			wantErr: true,
		},
		{
			name:    "another project",
			method:  &projects.Project{ID: 999, Name: "Another Project", APIKey: "api-key"},
			want:    prj,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, projects.ProjectKey, tt.method)

			got, err := projects.ProjectFromContext(ctx)

			if err != nil && !tt.wantErr {
				t.Errorf("ProjectFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("ProjectFromContext() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
