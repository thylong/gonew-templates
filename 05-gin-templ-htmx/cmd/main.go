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

	"github.com/spf13/cobra"
	"github.com/thylong/go-templates/05-gin-templ-htmx/internal/server"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	httpPort     string
	httpTimeout  int64
	loggingLevel string
	production   bool
)

func init() {
	rootCmd.AddCommand(versionCmd, runCmd)

	runCmd.Flags().StringVarP(&httpPort, "port", "p", "8080", "HTTP port to listen on")
	runCmd.Flags().StringVarP(&loggingLevel, "logging_level", "l", "info", "The app logging level")
	runCmd.Flags().Int64VarP(&httpTimeout, "timeout", "t", 500, "HTTP request timeout in milliseconds")
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

		r := server.CreateApp(production, httpTimeout)

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
