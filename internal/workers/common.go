package workers

import (
	"encoding/base64"
	"encoding/json"

	"github.com/fylerx/fyler/internal/jobs"
	"github.com/fylerx/fyler/internal/tasks"
)

func FetchTaskFromQueue(repo tasks.Repository, args ...interface{}) (*tasks.Task, error) {
	job := &jobs.Job{}

	rawDecodedText, err := base64.StdEncoding.DecodeString(args[0].(string))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawDecodedText, job)
	if err != nil {
		return nil, err
	}

	return repo.GetByID(job.TaskID)
}
