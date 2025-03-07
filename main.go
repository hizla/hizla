package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"git.gensokyo.uk/security/fortify/command"
	"github.com/hizla/hizla/hst"
	"github.com/hizla/hizla/internal"
	"github.com/hizla/hizla/internal/config"
	"github.com/hizla/hizla/internal/serve"
)

//go:embed LICENSE
var license string

func main() {
	log.SetFlags(0)
	log.SetPrefix("hizla: ")

	var flagVerbose bool

	c := command.New(os.Stderr, log.Printf, "hizla", func(args []string) error {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		return nil
	}).Flag(&flagVerbose, "v", command.BoolFlag(false), "Verbose output")

	c.Command("version", "Show hizla version", func([]string) error {
		if v, ok := internal.Check(internal.Version); ok {
			fmt.Println(v)
		} else {
			fmt.Println("impure")
		}
		return nil
	}).Command("license", "Show full license text", func([]string) error {
		fmt.Println(license)
		return nil
	}).Command("help", "Show this help message", func([]string) error {
		c.PrintHelp()
		return nil
	})

	{
		serveConfigPaths := [...]string{config.FromJSON: "", config.FromTOML: ""}

		c.NewCommand("serve", command.UsageInternal, func(args []string) error {
			var whence int
			for w, p := range serveConfigPaths {
				if p != "" {
					// guard against multiple paths
					if whence != config.FromEnviron {
						return errors.New("cannot load from multiple configuration types")
					}

					whence = w
				}
			}

			var r io.Reader
			if whence != config.FromEnviron {
				if f, err := os.Open(serveConfigPaths[whence]); err != nil {
					return fmt.Errorf("cannot open %q: %v", serveConfigPaths[whence], err)
				} else {
					r = f
				}
				log.Printf("loaded configuration from %q", serveConfigPaths[whence])
			}

			cfg := new(hst.ServeAPI)
			if err := config.Load(cfg, whence, r); err != nil {
				if errors.Is(err, hst.ErrUnsetAddress) {
					cfg.Address = defaultAPIListenAddr
					log.Printf("HIZLA_API_LISTEN_ADDRESS is unset, defaulting to %q", defaultAPIListenAddr)
				} else {
					return fmt.Errorf("cannot load config: %v", err)
				}
			}

			listenErr := make(chan error, 2)
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			if len(args) == 0 {
				startAPI(ctx, cfg, listenErr)
			} else {
				switch args[0] {
				case "api":
					startAPI(ctx, cfg, listenErr)
				default:
					return fmt.Errorf("%q is not a valid command", "serve "+args[0])
				}
			}

			if err := <-listenErr; err != nil {
				return err
			}
			if f := shutdownCancel.Load(); f != nil {
				(*f)()
			}
			return nil
		}).
			Flag(&serveConfigPaths[config.FromJSON], "j", command.StringFlag(""), "Path to JSON config file").
			Flag(&serveConfigPaths[config.FromTOML], "t", command.StringFlag(""), "Path to TOML config file")
	}

	c.MustParse(os.Args[1:], func(err error) {
		if err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

const (
	shutdownTimeout      = 5 * time.Second
	defaultAPIListenAddr = ":3000"
)

func startAPI(ctx context.Context, config *hst.ServeAPI, listenErr chan error) {
	shutdown := serve.StartAPI(config, listenErr)

	go func() {
		<-ctx.Done()
		log.Print("shutting down API server")
		if err := shutdown(serveShutdown()); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("cannot shut down API server: %v", err)
		}
	}()
}

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
