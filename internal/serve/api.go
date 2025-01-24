package serve

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hizla/hizla/hst"
	"github.com/hizla/hizla/internal"
	"github.com/hizla/hizla/internal/serve/handler"
)

var (
	version     string
	apiBaseResp []byte
)

func init() {
	if v, ok := internal.Check(internal.Version); !ok {
		version = "impure"
	} else {
		version = v
	}
}

func StartAPI(config *hst.ServeAPI, listenErr chan error) func(ctx context.Context) error {
	// pre-generate version response
	if b, err := json.Marshal(&hst.Version{Version: version}); err != nil {
		log.Fatalf("cannot serialise version response: %v", err)
	} else {
		apiBaseResp = b
	}

	app := fiber.New()
	app.Use(logger.New())

	v1 := app.Group("/api/v1")
	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.
			Type("json", "utf-8").
			Send(apiBaseResp)
	})

	app.Use(handler.CatchAll)
	go func() { listenErr <- app.Listen(config.Address) }()
	return app.ShutdownWithContext
}
