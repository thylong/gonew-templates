package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/03-k8s-fiber-sqlc/internal/server"
)

var (
	port         string
	httpTimeout  int64
	loggingLevel string
)

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&port, "port", "p", ":8080", "HTTP port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
	runCmd.Flags().Int64VarP(&httpTimeout, "timeout", "t", 500, "HTTP request timeout in milliseconds")
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

		// create Fiber app
		app := server.CreateApp(httpTimeout, loggingLevel)
		err := app.App.Listen("0.0.0.0" + port)
		if err != nil {
			log.Fatalf("fiber server failed to start: %v\n", err)
		}
	},
}
