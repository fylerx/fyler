package operations

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/fylerx/fyler/internal/conversions"
	"github.com/fylerx/fyler/internal/storage"
	"github.com/fylerx/fyler/internal/tasks"
)

type Operation struct {
	TaskID         uint64
	DownloadTime   int
	UploadTime     int
	FileSize       int64
	ResultPath     string
	SourceFilePath string
	TimeSpent      int
	Config         storage.Config
	Storage        storage.IClientS3
	Context        context.Context
}

func New(ctx context.Context, task *tasks.Task) (*Operation, error) {
	cfg := task.Project.Storage.Config()

	storage, err := storage.NewS3(cfg)
	if err != nil {
		return nil, err
	}

	return &Operation{
		TaskID:         task.ID,
		SourceFilePath: task.FilePath,
		Config:         cfg,
		Storage:        storage,
		Context:        ctx,
	}, nil
}

func (ops *Operation) UploadObject(file *os.File) error {
	startUpload := time.Now()
	location, err := ops.Storage.UploadObject(ops.Context, ops.Config.Bucket, file.Name(), file)
	if err != nil {
		return err
	}

	u, err := url.Parse(location)
	if err != nil {
		return err
	}
	ops.UploadTime = int(time.Since(startUpload).Seconds())
	ops.ResultPath = u.Path

	return err
}

func (ops *Operation) DownloadObject() (*os.File, error) {
	file, err := ioutil.TempFile("/tmp", "fylerx_")
	if err != nil {
		return file, err
	}

	startDownload := time.Now()
	err = ops.Storage.DownloadObject(ops.Context, ops.Config.Bucket, ops.SourceFilePath, file)
	ops.DownloadTime = int(time.Since(startDownload).Seconds())

	return file, err
}

func (ops *Operation) Conversion() *conversions.Conversion {
	return &conversions.Conversion{
		TaskID:       ops.TaskID,
		DownloadTime: ops.DownloadTime,
		UploadTime:   ops.UploadTime,
		FileSize:     ops.FileSize,
		ResultPath:   ops.ResultPath,
		TimeSpent:    ops.TimeSpent,
	}
}
