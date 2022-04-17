package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fylerx/fyler/internal"
)

func main() {
	worker := &internal.Worker{}
	if err := worker.Setup(); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("🚂 Starting worker...")
	go func() {
		if err := worker.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	log.Println("Shutting worker... Reason:", sig)

	if err := worker.Shutdown(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Worker gracefully stopped")
}
