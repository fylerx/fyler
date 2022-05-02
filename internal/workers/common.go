package workers

import (
	"context"
	"encoding/base64"
	"encoding/json"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/fylerx/fyler/internal/jobs"
	"github.com/fylerx/fyler/internal/tasks"
)

func FetchTaskFromQueue(ctx context.Context, repo tasks.Repository, arg string) (*tasks.Task, error) {
	job := &jobs.Job{}

	rawDecodedText, err := base64.StdEncoding.DecodeString(arg)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawDecodedText, job)
	if err != nil {
		return nil, err
	}

	help := worker.HelperFor(ctx)
	err = repo.SetProgressStatus(&tasks.Task{ID: job.TaskID}, help.Jid())
	if err != nil {
		return nil, err
	}

	return repo.GetByID(job.TaskID)
}
