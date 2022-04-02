package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/projects"
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

func (d *Worker) convertToPDF(ctx context.Context, args ...interface{}) error {
	help := worker.HelperFor(ctx)
	log.Println(args...)

	rawDecodedText, err := base64.StdEncoding.DecodeString(args[0].(string))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decoded text: %s\n", rawDecodedText)

	taskInput := tasks.Task{}
	err = json.Unmarshal(rawDecodedText, &taskInput)
	if err != nil {
		log.Printf("error %v\n", err)
	}
	log.Printf("W %v\n", taskInput)

	projectRepo := projects.InitRepo(d.repo)
	pj, err := projectRepo.GetByID(taskInput.ProjectID)
	if err != nil {
		log.Printf("error %v\n", err)
	}

	s3 := pj.Storage.Config()
	log.Printf("Ww %v\n", pj.Storage)

	log.Printf("==S3== %v\n", s3)
	session, err := storage.New(s3)
	// session, err := storage.New(storage.Config{
	// 	AccessKeyID:     s3.AccessKeyID.Raw.(string),
	// 	SecretAccessKey: s3.SecretAccessKey.Raw.(string),
	// 	Bucket:          s3.Bucket,
	// 	Endpoint:        s3.Endpoint,
	// 	Region:          s3.Region,
	// 	DisableSSL:      s3.DisableSSL,
	// })
	if err != nil {
		log.Printf("error %v\n", err)
	}

	cl := storage.NewS3(session, time.Second*5)

	file, err := os.OpenFile("testfile.pdf", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	res, err := cl.UploadObject(ctx, s3.Bucket, "testfile.pdf", file)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Working on job %s\n", help.Jid())
	log.Printf("Working on job %s\n", help.JobType())
	log.Printf("Working on res %s\n", res)
	return nil
}

func main() {

	cfg := &config.Config{}
	_, err := config.Read(constants.AppName, config.Defaults, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	aes, err := algorithms.NewAES256GCM([]byte(cfg.CRYPTO.Passphrase))
	if err != nil {
		log.Fatalf("failed to initialize crypto algorithm: %v\n", err)
	}
	gormcrypto.Init(aes, serialization.NewJSON())

	db, err := orm.Init(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// ctx, cancel := context.WithCancel(context.Background())
	mgr := worker.NewManager()

	mgr.Concurrency = 2

	mgr.Labels = append(mgr.Labels, "worker")

	mgr.ProcessStrictPriorityQueues("urgent", "high", "medium", "low")

	wr := &Worker{config: cfg, repo: db, wm: mgr}
	mgr.Register("doc_to_pdf", wr.convertToPDF)

	log.Println("Starting worker...")
	go func() {
		if err := mgr.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	log.Println("Shutting worker... Reason:", sig)

	// if err := dispatcher.Shutdown(); err != nil {
	// 	log.Fatal(err.Error())
	// }

	log.Println("Worker gracefully stopped")
}
