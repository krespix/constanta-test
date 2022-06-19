package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/krespix/constanta-test/internal/app/server"
	reqs "github.com/krespix/constanta-test/internal/services/request"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-c
		log.Printf("system call: %v", osCall)
		cancel()
	}()

	reqService := reqs.New(time.Second*1, 4)
	srv := server.New(":7000", reqService, 100)
	err := srv.Start(ctx)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
