package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	faktory "github.com/contribsys/faktory/client"
	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/handlers"
	"github.com/fylerx/fyler/internal/middleware"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/projects"
	"github.com/fylerx/fyler/internal/tasks"
	"github.com/gorilla/mux"
	gormcrypto "github.com/pkasila/gorm-crypto"
	"github.com/pkasila/gorm-crypto/algorithms"
	"github.com/pkasila/gorm-crypto/serialization"
	"gorm.io/gorm"
)

type Dispatcher struct {
	config *config.Config
	router *mux.Router
	repo   *gorm.DB
	jm     *faktory.Client
}

func (d *Dispatcher) Setup() error {
	cfg := &config.Config{}
	_, err := config.Read(constants.AppName, config.Defaults, cfg)
	if err != nil {
		return fmt.Errorf("[startup] can't read config, err: %w", err)
	}
	d.config = cfg

	aes, err := algorithms.NewAES256GCM([]byte(cfg.CRYPTO.Passphrase))
	if err != nil {
		log.Fatalf("failed to initialize crypto algorithm: %v\n", err)
	}
	gormcrypto.Init(aes, serialization.NewJSON())

	db, err := orm.Init(cfg)
	if err != nil {
		return fmt.Errorf("failed to init psql connection: %w", err)
	}
	d.repo = db

	client, err := faktory.Open()
	if err != nil {
		return fmt.Errorf("failed to init faktory connection: %w", err)
	}
	d.jm = client
	d.router = mux.NewRouter()

	d.initializeRoutes()

	return nil
}

func (d *Dispatcher) initializeRoutes() {
	tasksRepo := tasks.InitRepo(d.repo)
	handlers := &handlers.TasksHandler{TasksRepo: tasksRepo, JM: d.jm}

	d.router.HandleFunc("/api/tasks/", handlers.Index)
	d.router.HandleFunc("/api/tasks", handlers.Create).Methods("POST")

	projectsRepo := projects.InitRepo(d.repo)

	pmw := middleware.ProjectMiddleware{Projects: projectsRepo}
	d.router.Use(pmw.Middleware)
}

func (d *Dispatcher) ListenAndServe() error {
	port := strconv.Itoa(d.config.HTTP.Port)
	log.Printf("Connect to http://localhost:%s", port)

	return http.ListenAndServe(":"+port, d.router)
}

func (d *Dispatcher) Shutdown() error {
	dbConn, err := d.repo.DB()
	if err != nil {
		return fmt.Errorf("failed to get psql connection: %w", err)
	}

	fmt.Println("[PostrgeSQL] Closing connection...")
	if err = dbConn.Close(); err != nil {
		return fmt.Errorf("failed to close psql connection: %w", err)
	}

	fmt.Println("[Faktory] Closing connection...")
	if err = d.jm.Close(); err != nil {
		return fmt.Errorf("failed to close faktory connection: %w", err)
	}

	return nil
}
