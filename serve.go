package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hizla/hizla/hst"
	"github.com/hizla/hizla/internal/serve"
)

func doServe(args []string) {
	set := flag.NewFlagSet("serve", flag.ExitOnError)
	var apiListenAddr string
	set.StringVar(&apiListenAddr, "l", ":3000", "Print instance id")

	// Ignore errors; set is set for ExitOnError.
	_ = set.Parse(args[1:])

	listenErr := make(chan error, 2)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if len(set.Args()) == 0 {
		startAPI(ctx, &hst.ServeAPI{Address: apiListenAddr}, listenErr)
	} else {
		switch set.Args()[0] {
		case "api":
			startAPI(ctx, &hst.ServeAPI{Address: apiListenAddr}, listenErr)
		default:
			log.Fatal("invalid argument")
		}
	}

	if err := <-listenErr; err != nil {
		log.Fatalf("cannot listen: %v", err)
	}
	if f := shutdownCancel.Load(); f != nil {
		(*f)()
	}
}

func startAPI(ctx context.Context, config *hst.ServeAPI, listenErr chan error) {
	shutdown := serve.StartAPI(config, listenErr)

	go func() {
		<-ctx.Done()
		log.Print("shutting down API server")
		if err := shutdown(serveShutdown()); err != nil {
			log.Printf("cannot shut down API server: %v", err)
		}
	}()
}

const shutdownTimeout = 5 * time.Second

var (
	shutdownCtx    context.Context
	shutdownCancel atomic.Pointer[func()]
	shutdownOnce   sync.Once
)

func serveShutdown() context.Context {
	shutdownOnce.Do(func() {
		var cancel func()
		shutdownCtx, cancel = context.WithTimeout(context.Background(), shutdownTimeout)
		shutdownCancel.Store(&cancel)
	})
	return shutdownCtx
}
