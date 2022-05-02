package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/operations"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/tasks"
	"github.com/fylerx/fyler/internal/workers"
	"github.com/go-resty/resty/v2"
	gormcrypto "github.com/pkasila/gorm-crypto"
	"github.com/pkasila/gorm-crypto/algorithms"
	"github.com/pkasila/gorm-crypto/serialization"
)

type Worker struct {
	config *config.Config
	wm     *worker.Manager
	tasks  tasks.Repository
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
	task, err := workers.FetchTaskFromQueue(ctx, w.tasks, args[0].(string))
	if err != nil {
		log.Printf("error %v\n", err)
		return err
	}

	operation, err := operations.New(ctx, task)
	if err != nil {
		log.Printf("error %v\n", err)
		w.tasks.Failed(task, err)
		return err
	}

	file, err := operation.DownloadObject()
	defer os.Remove(file.Name())
	if err != nil {
		log.Printf("error %v\n", err)
		w.tasks.Failed(task, err)
		return err
	}

	// Convert file
	startTimeSpent := time.Now()
	FileOut := fmt.Sprintf("converted_%s", file.Name())
	out, err := exec.Command("unoconv", "-o", FileOut, "-f", "pdf", file.Name()).Output()
	if err != nil {
		w.tasks.Failed(task, err)
		log.Printf("error %v\n", err)
		return err
	}
	fmt.Println("Command Successfully Executed", out)

	fi, err := file.Stat()
	if err != nil {
		w.tasks.Failed(task, err)
		log.Printf("error %v\n", err)
		return err
	}
	operation.TimeSpent = int(time.Since(startTimeSpent).Seconds())
	operation.FileSize = fi.Size()
	operation.ResultPath = fi.Name()
	// -

	// Uploading file
	err = operation.UploadObject(file)
	if err != nil {
		log.Printf("error %v\n", err)
		w.tasks.Failed(task, err)
		return err
	}
	task.Conversion = operation.Conversion()
	err = w.tasks.SetSuccessStatus(task)
	if err != nil {
		w.tasks.Failed(task, err)
		log.Printf("error %v\n", err)
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(task.Conversion).
		Post(task.Project.CallbackURL)
	fmt.Printf("resp %v", resp)
	return nil
}
