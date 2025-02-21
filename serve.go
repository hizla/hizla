package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hizla/hizla/hst"
	"github.com/hizla/hizla/internal/config"
	"github.com/hizla/hizla/internal/serve"
)

const defaultAPIListenAddr = ":3000"

func doServe(args []string) {
	set := flag.NewFlagSet("serve", flag.ExitOnError)
	configPaths := [...]*string{
		config.FromJSON: set.String("j", "", "Path to JSON config file"),
		config.FromTOML: set.String("t", "", "Path to TOML config file"),
	}

	// Ignore errors; set is set for ExitOnError.
	_ = set.Parse(args[1:])

	var whence int
	for w, p := range configPaths {
		if p != nil && *p != "" {
			// guard against multiple paths
			if whence != config.FromEnviron {
				log.Fatal("cannot load from multiple configuration types")
			}

			whence = w
		}
	}

	var r io.Reader
	if whence != config.FromEnviron {
		if f, err := os.Open(*configPaths[whence]); err != nil {
			log.Fatalf("cannot open %q: %v", *configPaths[whence], err)
		} else {
			r = f
		}
	}

	cfg := new(hst.ServeAPI)
	if err := config.Load(cfg, whence, r); err != nil {
		if errors.Is(err, hst.ErrUnsetAddress) {
			cfg.Address = defaultAPIListenAddr
			log.Printf("HIZLA_API_LISTEN_ADDRESS is unset, defaulting to %q", defaultAPIListenAddr)
		} else {
			log.Fatalf("cannot load config: %v", err)
		}
	}

	listenErr := make(chan error, 2)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if len(set.Args()) == 0 {
		startAPI(ctx, cfg, listenErr)
	} else {
		switch set.Args()[0] {
		case "api":
			startAPI(ctx, cfg, listenErr)
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
