package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/db"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/middleware"
	eventpb "github.com/thylong/go-templates/06-grpc-sqlc/pkg/proto"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	production   bool
	port         string
	httpTimeout  int64
	databaseURI  string
	loggingLevel string
)

const version = "v0.1.1"

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&port, "port", "p", "50051", "gRPC port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
	runCmd.Flags().BoolVarP(&production, "production", "g", false, "enable production settings (logging fmt, prefork, etc)")
	runCmd.Flags().Int64VarP(&httpTimeout, "timeout", "t", 500, "HTTP request timeout in milliseconds")
	runCmd.Flags().StringVarP(&databaseURI, "database", "c", "postgres://admin:secret@db:5432/postgres?sslmode=disable", "Postgresql database URI, default to local Docker env")
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
		// Set up zerolog
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

		flag.Parse()

		pool, err := pgxpool.New(context.Background(), databaseURI)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
		defer pool.Close()

		// Create db.Queries using the pool
		queries := db.New(pool)

		// Set up gRPC server options with middleware
		opts := []logging.Option{
			logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			// Add any other option (check functions starting with logging.With).
		}

		// You can now create a server with logging instrumentation that e.g. logs when the unary or stream call is started or finished.
		app := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				middleware.DefaultTimeoutUnaryInterceptor(time.Duration(httpTimeout)*time.Millisecond),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					log.Error().Msgf("panic recovered: %v", p)
					return status.Errorf(codes.Internal, "unexpected error occurred")
				})),
				logging.UnaryServerInterceptor(InterceptorLogger(log.Logger), opts...),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(InterceptorLogger(log.Logger), opts...),
			),
		)
		healthServer := health.NewServer()

		grpc_health_v1.RegisterHealthServer(app, healthServer)
		eventpb.RegisterEventServiceServer(app, server.NewEventServiceServer(queries))
		// Set the health status of the server
		healthServer.SetServingStatus("EventService", grpc_health_v1.HealthCheckResponse_SERVING)
		reflection.Register(app)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Sprintf("failed to start listener: %v", err))
		}

		// Handle graceful shutdown
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			fmt.Printf("gRPC server running on port %s\n", port)
			if err := app.Serve(listener); err != nil {
				log.Fatal().Err(err).Msg("failed to serve")
			}
		}()

		<-stop
		fmt.Println("Shutting down gRPC server...")
		app.GracefulStop()
		fmt.Println("gRPC server stopped")
	},
}

// InterceptorLogger adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}
