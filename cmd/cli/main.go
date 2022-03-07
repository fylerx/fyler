package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/projects"
	"github.com/urfave/cli/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := &config.Config{}
	_, err := config.Read(constants.AppName, config.Defaults, cfg)
	if err != nil {
		log.Fatalf("[startup] can't read config, err: %v\n", err)
	}

	db, err := orm.Init(cfg)
	if err != nil {
		log.Fatalf("failed to init psql connection: %v\n", err)
	}

	pj := projects.InitRepo(db)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"new"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					println(c.Args().First())
					data := &projects.Project{Name: c.Args().First()}
					p, err := pj.Create(data)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("API-KEY:", p.APIKey)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "complete a task on the list",
				Action: func(c *cli.Context) error {
					fmt.Println(pj.GetAll())
					fmt.Println("completed task: ", c.Args().First())
					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
