package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	gormcrypto "github.com/pkasila/gorm-crypto"
	"github.com/pkasila/gorm-crypto/algorithms"
	"github.com/pkasila/gorm-crypto/serialization"
	"gorm.io/gorm"
)

type Worker struct {
	config *config.Config
	repo   *gorm.DB
	wm     *worker.Manager
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
	w.repo = db

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
	dbConn, err := w.repo.DB()
	if err != nil {
		return fmt.Errorf("failed to get psql connection: %w", err)
	}

	fmt.Println("[PostrgeSQL] Closing connection...")
	if err = dbConn.Close(); err != nil {
		return fmt.Errorf("failed to close psql connection: %w", err)
	}

	fmt.Println("[Faktory Worker] Closing connection...")
	w.wm.Quiet()

	return nil
}

func (w *Worker) convertToPDF(ctx context.Context, args ...interface{}) error {
	// Get job
	help := worker.HelperFor(ctx)
	rawDecodedText, err := base64.StdEncoding.DecodeString(args[0].(string))
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("Decoded text: %s\n", rawDecodedText)

	taskJob := &tasks.Task{}
	err = json.Unmarshal(rawDecodedText, taskJob)
	if err != nil {
		log.Printf("error %v\n", err)
		return err
	}

	taskRepo := tasks.InitRepo(w.repo)
	task, err := taskRepo.GetByID(taskJob.ID)
	if err != nil {
		log.Printf("error %v\n", err)
		return err
	}
	taskRepo.Progress(task)

	s3 := task.Project.Storage.Config()
	session, err := storage.New(s3)
	if err != nil {
		log.Printf("error %v\n", err)
	}

	clientS3 := storage.NewS3(session, time.Second*5)
	conv := &conversions.Conversion{
		TaskID: task.ID,
		JobID:  help.Jid(),
	}

	// Downloading file
	file, err := ioutil.TempFile("/tmp", "fylerx_")
	if err != nil {
		log.Printf("error %v\n", err)
		taskRepo.Failed(task, err)
		return err
	}
	defer os.Remove(file.Name())

	startDownload := time.Now()
	err = clientS3.DownloadObject(ctx, s3.Bucket, task.FilePath, file)
	if err != nil {
		log.Printf("error %v\n", err)
		taskRepo.Failed(task, err)
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
		taskRepo.Failed(task, err)
		log.Printf("error %v\n", err)
	}
	conv.UploadTime = int(time.Since(startUpload).Seconds())
	task.Conversion = conv

	err = taskRepo.Success(task)
	if err != nil {
		taskRepo.Failed(task, err)
		log.Printf("error %v\n", err)
	}

	log.Printf("Working on job %s\n", help.Jid())
	log.Printf("Working on job %s\n", help.JobType())
	log.Printf("Working on res %s\n", res)
	log.Printf("Working on conv %v\n", conv)
	return nil
}
