package main

import (
	"context"
	"log"
	"time"

	"github.com/krespix/constanta-test/internal/app/server"
	reqs "github.com/krespix/constanta-test/internal/services/request"
)

func main() {
	ctx := context.Background()
	reqService := reqs.New(time.Second*1, 4)
	srv := server.New(":7000", reqService, 100)
	err := srv.Start(ctx)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
