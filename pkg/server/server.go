package server

import (
	"context"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"github.com/williamnoble/kube-botany/pkg/render"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

// Server represents the HTTP httpServer for the plant application
type Server struct {
	staticDir    string                        // Directory for static assets
	templatesDir string                        // Directory for HTML templates
	templates    map[string]*template.Template // Parsed HTML templates

	Logger    *slog.Logger // Logger for httpServer logs
	startTime time.Time    // Time when the httpServer started

	mu     sync.Mutex     // Mutex for thread-safe access to plants
	plants []*plant.Plant // Collection of plants managed by the httpServer

	renderer *render.ASCIIRenderer // Renderer for ASCII art

	httpServer *http.Server
}

// NewServer creates a new httpServer instance with the given plants
// It initializes the logger, renderer, templates, and other httpServer components
func NewServer(plants []*plant.Plant) *Server {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	s := &Server{
		plants:       plants,
		Logger:       logger,
		startTime:    time.Now(),
		renderer:     render.NewASCIIRenderer(),
		templates:    make(map[string]*template.Template),
		staticDir:    "cmd/api/static",
		templatesDir: "cmd/api/templates",
	}
	s.ParseTemplates()

	return s
}

// Start starts the HTTP httpServer on the specified port
// It sets up the routes, adds the request logger middleware, and starts listening for requests
func (s *Server) Start(port int) error {
	mux := s.Routes()
	addr := fmt.Sprintf(":%d", port)

	s.Logger.Info("starting httpServer", "addr", addr)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.requestLogger(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	s.Logger.Info("starting background tasks")
	go s.BackgroundTasks()

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		s.Logger.Error("http server not started")
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}
