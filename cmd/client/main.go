package main

import (
	"context"
	"gophkeeper/cmd/client/configuration"
	"gophkeeper/internal/app/client/cli"
	grpcclient "gophkeeper/internal/app/client/grpc"
	"gophkeeper/internal/app/client/storage"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	c := configuration.New()
	go handleSignals(cancel)
	gc, gcErr := grpcclient.New(c.GRPC_ADDR, time.Duration(c.GRPC_TIMEOUT))
	s := storage.New(c.STORAGE_PATH)
	var cl *cli.CLI
	if gcErr != nil {
		cl = cli.New(nil, s)
	} else {
		cl = cli.New(gc, s)
	}
	err := cl.Start(cancel)
	if err != nil {
		log.Err(err).Msg("error")
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	_, cancelt := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelt()
	select {
	case <-sigint:
		cancel()
		return
	}
}
