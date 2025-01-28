package serve

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hizla/hizla/hst"
	"github.com/hizla/hizla/internal"
	"github.com/hizla/hizla/internal/auth"
	"github.com/hizla/hizla/internal/database"
	"github.com/hizla/hizla/internal/serve/handler"
	"log"
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
	auth.InitStore()

	dbPool, err := database.NewDBPool(&database.DBConfig{
		Host:     config.DbHost,
		Port:     config.DbPort,
		User:     config.DbUser,
		Password: config.DbPassword,
		DBName:   config.DbName,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authHandler := auth.NewAuthHandler(dbPool)

	app.Post("/signin", authHandler.SignIn)
	app.Post("/signup", authHandler.SignUp)
	app.Post("/logout", authHandler.Logout)
	app.Get("/profile", authHandler.Profile)

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
