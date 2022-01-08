package internal

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	faktory "github.com/contribsys/faktory/client"
	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
	AppName = "dispatcher"
)

type Dispatcher struct {
	config *config.Config
	router *mux.Router
	repo   *gorm.DB
	jm     *faktory.Client
}

func (d *Dispatcher) Setup() error {
	cfg := &config.Config{}
	_, err := config.Read(AppName, config.Defaults, cfg)
	if err != nil {
		return fmt.Errorf("[startup] can't read config, err: %w", err)
	}
	d.config = cfg

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
	return nil
}

func (d *Dispatcher) ListenAndServe() error {
	port := strconv.Itoa(d.config.HTTP.Port)
	log.Printf("🚀 connect to http://localhost:%s", port)

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

	//TODO: Close faktory
	return nil
}
