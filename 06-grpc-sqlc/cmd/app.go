package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	_ "net/http/pprof"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/db"
	eventpb "github.com/thylong/go-templates/06-grpc-sqlc/pkg/proto"
	"github.com/thylong/go-templates/06-grpc-sqlc/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		flag.Parse()

		pool, err := pgxpool.New(context.Background(), databaseURI)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
		defer pool.Close()

		// Create db.Queries using the pool
		queries := db.New(pool)

		app := grpc.NewServer()
		eventpb.RegisterEventServiceServer(app, server.NewEventServiceServer(queries))
		reflection.Register(app)

		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			panic(fmt.Sprintf("failed to start listener: %v", err))
		}

		fmt.Printf("gRPC server running on port %s\n", port)
		if err := app.Serve(listener); err != nil {
			panic(fmt.Sprintf("failed to serve: %v", err))
		}
	},
}
