package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/db"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/handlers"
	"github.com/thylong/go-templates/04-gin-sqlc/pkg/middlewares"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	httpPort     string
	httpTimeout  int64
	databaseURI  string
	loggingLevel string
	production   bool
)

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&httpPort, "port", "p", "8080", "HTTP port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
	runCmd.Flags().Int64VarP(&httpTimeout, "timeout", "t", 500, "HTTP request timeout in milliseconds")
	runCmd.Flags().StringVarP(&databaseURI, "database", "c", "postgres://admin:secret@db:5432/postgres?sslmode=disable", "Postgresql database URI, default to local Docker env")
	runCmd.Flags().BoolVarP(&production, "production", "g", false, "enable production settings (logging fmt, prefork, etc)")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "return current app version",
	Long:  `Return current application version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the app",
	Long:  `Run the application with given configuration (default with optional CLI flags overrides)`,
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()

		tracer.Start()
		defer tracer.Stop()

		conn, err := pgx.Connect(context.Background(), databaseURI)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		q := db.New(conn)
		r := setupRouter(q)

		// Use the tracer middleware with given service name.
		r.Use(gintrace.Middleware("app"))

		// Listen and Serve in 0.0.0.0:<port> (default 8000)
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", httpPort),
			Handler: r,
		}

		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("listen: %s\n", err)
				os.Exit(1)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscanll.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("Server Shutdown:", err)
			os.Exit(1)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		<-ctx.Done()
		fmt.Println("Server exiting")
	},
}

func setupRouter(queries *db.Queries) *gin.Engine {
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

	// Healthcheck endpoint for k8s probes
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api")
	router := api.Group("/auth")
	router.POST("/register", handlers.NewAuthHandler(queries).SignUpUser)

	return r
}
