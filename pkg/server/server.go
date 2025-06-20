package server

import (
	"context"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/render"
	"github.com/williamnoble/kube-botany/pkg/repository"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Server represents the HTTP server for the plant application
type Server struct {
	staticDir string                        // Directory for static assets
	templates map[string]*template.Template // Parsed HTML templates

	Logger    *slog.Logger // Logger for server logs
	startTime time.Time    // Time when the server started

	store    repository.PlantRepository // Repository for plants
	renderer *render.ASCIIRenderer      // Renderer for ASCII art

	httpServer *http.Server
}

// NewServer creates a new Server instance with the given plants
// It initialises the logger, renderer, templates, and other server components
func NewServer(inMemoryStore repository.PlantRepository) (*Server, error) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	s := &Server{
		Logger:    logger,
		startTime: time.Now(),
		templates: make(map[string]*template.Template),
		staticDir: "pkg/static",
		store:     inMemoryStore,
		renderer:  render.NewASCIIRenderer(),
	}
	s.ParseTemplates()

	return s, nil
}

// Start starts the HTTP server on the specified port
// It sets up the routes, adds the request logger middleware, and starts listening for requests
func (s *Server) Start(portStr string) error {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	mux := s.Routes()
	addr := fmt.Sprintf(":%d", port)

	s.Logger.With("component", "server").Info("server", "listening on", addr)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.requestLogger(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		s.Logger.With("component", "server").Info("http server not started")
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}
