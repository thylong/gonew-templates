package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	production   bool
	port         string
	httpTimeout  int64
	loggingLevel string
)

const version = "v0.1.1"

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&port, "port", "p", "8080", "HTTP port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
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

		fmt.Println("喂，你好！")
	},
}
