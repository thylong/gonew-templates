package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/02-simple-k8s-fiber-app/internal/server"
)

var (
	production   bool
	port         string
	httpTimeout  int64
	loggingLevel string
	host         string
	enablePprof  bool
)

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&port, "port", "p", ":8080", "HTTP port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
	runCmd.Flags().BoolVarP(&enablePprof, "profile", "f", false, "enable profiling with Pprof")
	runCmd.Flags().BoolVarP(&production, "production", "g", false, "enable production settings (logging fmt, prefork, etc)")
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

		host := "localhost"
		if production {
			host = "0.0.0.0"
		}

		// if enablePprof {
		// 	go func() {
		// 		log.Printf("Starting pprof on %s:6060", host)
		// 		if err := http.ListenAndServe(host+":6060", nil); err != nil {
		// 			log.Fatalf("pprof server failed: %v\n", err)
		// 		}
		// 	}()
		// }

		// create Fiber app
		app := server.CreateApp(httpTimeout, loggingLevel, production)
		err := app.App.Listen(host + port)
		if err != nil {
			log.Fatalf("fiber server failed to start: %v\n", err)
		}
	},
}
