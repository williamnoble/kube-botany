package main

import (
	"context"
	"errors"
	"github.com/williamnoble/kube-botany/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	svr, err := server.NewServer(true)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := svr.Start(8090); err != nil && !errors.Is(err, http.ErrServerClosed) {
			svr.Logger.Error("server error", "error", err)
			panic(err)
		}
	}()

	svr.Logger.Info("server started successfully")

	// Block until we receive a signal
	<-quit
	svr.Logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		svr.Logger.Error("server forced to shutdown", "error", err)
		panic(err)
	}

	svr.Logger.Info("server exited gracefully")
}
