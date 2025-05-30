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
			svr.Logger.With("component", "server").Error("server error", "error", err)
			panic(err)
		}
	}()

	<-quit
	shutdownTimeout := 10 * time.Second
	svr.Logger.With("component", "server").Info("starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		svr.Logger.With("component", "server").Error("graceful shutdown failed", err)
		panic(err)
	}

	svr.Logger.With("component", "server").Info("graceful shutdown complete")
}
