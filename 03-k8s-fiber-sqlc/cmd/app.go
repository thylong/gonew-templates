package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/internal/server"
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

	runCmd.Flags().StringVarP(&httpPort, "port", "p", ":8080", "HTTP port to listen on")
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
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		flag.Parse()

		conn, err := pgx.Connect(context.Background(), databaseURI)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		app := server.CreateApp(httpTimeout, loggingLevel, production, conn)
		err = app.App.Listen("0.0.0.0" + httpPort)
		if err != nil {
			log.Fatalf("fiber server failed to start: %v\n", err)
		}
	},
}
