package server

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/internal/core"
)

func CreateApp(httpTimeout int64, loggingLevel string) *core.App {
	fiberApp := fiber.New(fiber.Config{
		Prefork:               true,
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		},
	})

	// init store logger
	if loggingLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Setting global log level to debug")
	}

	// Middlewares
	fiberApp.Use(recover.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(compress.New())

	// Catch all handler
	fiberApp.Use(timeout.NewWithContext(
		func(c *fiber.Ctx) error {
			return c.SendStatus(500)
		},
		time.Duration(httpTimeout)*time.Millisecond),
	)

	app := &core.App{App: fiberApp}

	return app
}
