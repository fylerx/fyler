package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fylerx/fyler/internal"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	dispatcher := &internal.Dispatcher{}
	if err := dispatcher.Setup(); err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := dispatcher.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	if err := dispatcher.Shutdown(); err != nil {
		log.Fatal(err.Error())
	}
}
