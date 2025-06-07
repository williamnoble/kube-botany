package main

import (
	"context"
	"errors"
	"github.com/williamnoble/kube-botany/pkg/config"
	"github.com/williamnoble/kube-botany/pkg/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := config.NewFromEnvironment()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	svr, err := server.NewServer(true)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := svr.Start(c.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			svr.Logger.With("component", "server").Error("server error", "error", err)
			panic(err)
		}
	}()

	go svr.BackgroundTasks(ctx)

	<-ctx.Done()
	shutdownTimeout := 10 * time.Second
	svr.Logger.With("component", "server").Info("starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		svr.Logger.With("component", "server").Error("graceful shutdown failed", "error", err)
		panic(err)
	}

	svr.Logger.With("component", "server").Info("graceful shutdown complete")
}
