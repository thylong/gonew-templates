package server

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/thylong/go-templates/04-gin-sqlc/api"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/db"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/middlewares"
)

func CreateApp(production bool, httpTimeout int64, queries *db.Queries) *gin.Engine {
	r := gin.New()

	/*
	   database setup:
	   - Initialize connection to the database
	   - Start session and store in gin context

	   Gin setup:
	   - Chain global middlewares (logger, timeout, recovery, etc)
	   - Setup configuration (dev vs prod, log format, gracefull shutdown)
	   - Create healthcheck route /healthz
	   - Setup configuration (dev vs prod, log format)
	*/

	// Disable Console Color, switch to gin ReleaseMode, configure logger.
	if production {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)

		// Logger middleware will:
		//   - Logs all requests, like a combined access and error log.
		//   - Logs to stdout.
		//   - Logs using JSON format.
		r.Use(logger.SetLogger(logger.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().Logger()
		})))

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		// Logger middleware will:
		//   - Logs all requests, like a combined access and error log.
		//   - Logs to stdout.
		r.Use(logger.SetLogger())
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// timeout middleware will return 408 if timeout is reached.
	r.Use(middlewares.TimeoutMiddleware(httpTimeout))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	api.SetupRoutes(r, queries)

	return r
}
