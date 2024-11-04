package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"dashboardk/internal/config"
	"dashboardk/internal/server"

	autorefresh "github.com/lavigneer/browser-autorefresh"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func gracefulShutdown(apiServer *http.Server, done chan bool, app *config.Application) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	app.Logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		app.Logger.Error("Server forced to shutdown with error", "error", err)
	}

	app.Logger.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func createClientSet() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	appEnv := os.Getenv("APP_ENV")
	var pageReloader *autorefresh.PageReloader
	if appEnv == "local" {
		var err error
		pageReloader, err = autorefresh.New(nil, "/__dev__/reload", 100)
		if err != nil {
			logger.Error("could not set up autorefresh page reloader", "error", err)
		}
	}
	app := &config.Application{
		Logger:       logger,
		ClientSet:    createClientSet(),
		Env:          os.Getenv("APP_ENV"),
		PageReloader: pageReloader,
	}
	server := server.NewServer(app)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done, app)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	app.Logger.Info("Graceful shutdown complete.")
}
