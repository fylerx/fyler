package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fylerx/fyler/internal/config"
	"github.com/fylerx/fyler/internal/constants"
	"github.com/fylerx/fyler/internal/orm"
	"github.com/fylerx/fyler/internal/projects"
	"github.com/fylerx/fyler/internal/storages"
	"github.com/mitchellh/mapstructure"
	gormcrypto "github.com/pkasila/gorm-crypto"
	"github.com/pkasila/gorm-crypto/algorithms"
	"github.com/pkasila/gorm-crypto/serialization"
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

	aes, err := algorithms.NewAES256GCM([]byte(cfg.CRYPTO.Passphrase))
	if err != nil {
		log.Fatalf("failed to initialize crypto algorithm: %v\n", err)
	}
	gormcrypto.Init(aes, serialization.NewJSON())

	projectRepo := projects.InitRepo(db)

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "new",
				Aliases: []string{"new"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					println(c.Args().First())
					data := &projects.Project{Name: c.Args().First()}
					p, err := projectRepo.Create(data)
					if err != nil {
						log.Fatal(err)
					}
					storage := storages.InitRepo(db)

					input := [6]string{
						"access_key",
						"secret_key",
						"bucket",
						"endpoint",
						"region",
						"disable_ssl",
					}
					res := make(map[string]interface{})
					for _, val := range input {
						res[val] = StringPrompt(val)
					}

					var s3 storages.Storage

					dc := &mapstructure.DecoderConfig{
						Result: &s3,
						DecodeHook: mapstructure.ComposeDecodeHookFunc(
							StringToBoolHookFunc,
							StringToCryptedHookFunc,
						)}
					ms, err := mapstructure.NewDecoder(dc)
					if err != nil {
						return err
					}

					err = ms.Decode(res)
					if err != nil {
						log.Fatal(err)
					}

					s3.ProjectID = p.ID
					err = storage.CreateStorage(&s3)
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
					projects, _ := projectRepo.GetAll()
					for _, pj := range projects {
						fmt.Printf("%d. %s\n", pj.ID, pj.Name)
					}
					return nil
				},
			},
			{
				Name:    "get",
				Aliases: []string{"get"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					u64, err := strconv.ParseUint(c.Args().First(), 10, 32)
					if err != nil {
						log.Fatal(err)
					}

					p, err := projectRepo.GetByID(uint32(u64))
					if err != nil {
						log.Fatal(err)
					}

					payload, _ := json.Marshal(p.Storage)

					fmt.Println("Storage:", string(payload))
					return nil
				},
			},
			{
				Name:    "add_storage",
				Aliases: []string{"add_storage"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					u64, err := strconv.ParseUint(c.Args().First(), 10, 32)
					if err != nil {
						log.Fatal(err)
					}

					p, err := projectRepo.GetByID(uint32(u64))
					if err != nil {
						log.Fatal(err)
					}

					storage := storages.InitRepo(db)

					input := [6]string{
						"access_key",
						"secret_key",
						"bucket",
						"endpoint",
						"region",
						"disable_ssl",
					}
					res := make(map[string]interface{})
					for _, val := range input {
						res[val] = StringPrompt(val)
					}

					var s3 storages.Storage

					dc := &mapstructure.DecoderConfig{
						Result: &s3,
						DecodeHook: mapstructure.ComposeDecodeHookFunc(
							StringToBoolHookFunc,
							StringToCryptedHookFunc,
						)}
					ms, err := mapstructure.NewDecoder(dc)
					if err != nil {
						return err
					}

					err = ms.Decode(res)
					if err != nil {
						log.Fatal(err)
					}

					s3.ProjectID = p.ID
					err = storage.CreateStorage(&s3)
					if err != nil {
						log.Fatal(err)
					}

					payload, _ := json.Marshal(s3)

					fmt.Println("Storage:", string(payload))
					return nil
				},
			},
			{
				Name:    "delete_storage",
				Aliases: []string{"delete_storage"},
				Usage:   "add a task to the list",
				Action: func(c *cli.Context) error {
					u64, err := strconv.ParseUint(c.Args().First(), 10, 32)
					if err != nil {
						log.Fatal(err)
					}

					p, err := projectRepo.GetByID(uint32(u64))
					if err != nil {
						log.Fatal(err)
					}

					storage := storages.InitRepo(db)
					err = storage.DeleteStorage(p.Storage.ID)
					if err != nil {
						log.Fatal(err)
					}
					payload, _ := json.Marshal(p.Storage)

					fmt.Println("Storage:", string(payload))
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

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func StringToBoolHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() == reflect.String && t.Kind() == reflect.Bool {
		val, _ := strconv.ParseBool(data.(string))
		return val, nil
	}

	return data, nil
}

func StringToCryptedHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	fmt.Println(t.Kind())
	if f.Kind() == reflect.String && t.Kind() == reflect.Struct {
		crypted := gormcrypto.EncryptedValue{Raw: data.(string)}

		return crypted, nil
	}

	return data, nil
}
