package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/conversions"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/storage"
	"github.com/fylerx/fyler/internal/tasks"
	"github.com/fylerx/fyler/internal/workers"
	gormcrypto "github.com/pkasila/gorm-crypto"
	"github.com/pkasila/gorm-crypto/algorithms"
	"github.com/pkasila/gorm-crypto/serialization"
)

type Worker struct {
	config *config.Config
	// repo   *gorm.DB
	wm    *worker.Manager
	tasks tasks.Repository
}

func (w *Worker) Setup() error {
	cfg := &config.Config{}
	_, err := config.Read(constants.AppName, config.Defaults, cfg)
	if err != nil {
		return fmt.Errorf("[startup] can't read config, err: %w", err)
	}
	w.config = cfg

	aes, err := algorithms.NewAES256GCM([]byte(cfg.CRYPTO.Passphrase))
	if err != nil {
		log.Fatalf("failed to initialize crypto algorithm: %v\n", err)
	}
	gormcrypto.Init(aes, serialization.NewJSON())

	db, err := orm.Init(cfg)
	if err != nil {
		return fmt.Errorf("failed to init psql connection: %w", err)
	}
	w.tasks = tasks.InitRepo(db)

	mgr := worker.NewManager()
	mgr.Concurrency = 2
	mgr.Labels = append(mgr.Labels, "worker")
	mgr.ProcessStrictPriorityQueues("urgent", "high", "medium", "low")
	w.wm = mgr
	mgr.Register("doc_to_pdf", w.convertToPDF)
	return nil
}

func (w *Worker) Run() error {
	return w.wm.Run()
}

func (w *Worker) Shutdown() error {
	fmt.Println("[PostrgeSQL] Closing connection...")
	if err := w.tasks.CloseDBConnection(); err != nil {
		return fmt.Errorf("failed to close psql connection: %w", err)
	}

	fmt.Println("[Faktory Worker] Closing connection...")
	w.wm.Quiet()

	return nil
}

func (w *Worker) convertToPDF(ctx context.Context, args ...interface{}) error {
	// Get job
	help := worker.HelperFor(ctx)
	task, err := workers.FetchTaskFromQueue(w.tasks, args)
	if err != nil {
		log.Printf("error %v\n", err)
		return err
	}

	err = w.tasks.SetProgressStatus(task)
	if err != nil {
		log.Printf("error %v\n", err)
		return err
	}

	conv := &conversions.Conversion{
		TaskID: task.ID,
		JobID:  help.Jid(),
	}

	s3 := task.Project.Storage.Config()
	session, err := storage.New(s3)
	if err != nil {
		log.Printf("error %v\n", err)
	}

	clientS3 := storage.NewS3(session, time.Second*5)

	// Downloading file
	file, err := ioutil.TempFile("/tmp", "fylerx_")
	if err != nil {
		log.Printf("error %v\n", err)
		w.tasks.Failed(task, err)
		return err
	}
	defer os.Remove(file.Name())

	startDownload := time.Now()
	err = clientS3.DownloadObject(ctx, s3.Bucket, task.FilePath, file)
	if err != nil {
		log.Printf("error %v\n", err)
		w.tasks.Failed(task, err)
		return err
	}
	conv.DownloadTime = int(time.Since(startDownload).Seconds())
	// -

	// Convert file
	startTimeSpent := time.Now()
	FileOut := fmt.Sprintf("converted_%s", file.Name())
	out, err := exec.Command("unoconv", "-o", FileOut, "-f", "pdf", file.Name()).Output()
	if err != nil {
		// taskRepo.Failed(taskInput.ID, err)
		log.Printf("error %v\n", err)
		// return err
	}
	fmt.Println("Command Successfully Executed", out)

	fi, err := file.Stat()
	if err != nil {
		// taskRepo.Failed(taskInput.ID, err)
		log.Printf("error %v\n", err)
		// return err
	}
	conv.TimeSpent = int(time.Since(startTimeSpent).Seconds())
	conv.FileSize = fi.Size()
	conv.ResultPath = fi.Name()
	// -

	// Uploading file
	startUpload := time.Now()
	res, err := clientS3.UploadObject(ctx, s3.Bucket, file.Name(), file)
	if err != nil {
		w.tasks.Failed(task, err)
		log.Printf("error %v\n", err)
	}
	conv.UploadTime = int(time.Since(startUpload).Seconds())
	task.Conversion = conv

	err = w.tasks.SetSuccessStatus(task)
	if err != nil {
		w.tasks.Failed(task, err)
		log.Printf("error %v\n", err)
	}

	log.Printf("Working on job %s\n", help.Jid())
	log.Printf("Working on job %s\n", help.JobType())
	log.Printf("Working on res %s\n", res)
	log.Printf("Working on conv %v\n", conv)
	return nil
}
