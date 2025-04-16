package server

import (
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

// Server represents the HTTP server for the plant application
type Server struct {
	staticDir    string                        // Directory for static assets
	templatesDir string                        // Directory for HTML templates
	templates    map[string]*template.Template // Parsed HTML templates

	logger    *slog.Logger // Logger for server logs
	startTime time.Time    // Time when the server started

	mu     sync.Mutex     // Mutex for thread-safe access to plants
	plants []*plant.Plant // Collection of plants managed by the server

	renderer *render.ASCIIRenderer // Renderer for ASCII art
}

// NewServer creates a new server instance with the given plants
// It initializes the logger, renderer, templates, and other server components
func NewServer(plants []*plant.Plant) *Server {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	s := &Server{
		plants:       plants,
		logger:       logger,
		startTime:    time.Now(),
		renderer:     render.NewASCIIRenderer(),
		templates:    make(map[string]*template.Template),
		staticDir:    "cmd/api/static",
		templatesDir: "cmd/api/templates",
	}
	s.ParseTemplates()

	return s
}

// Start starts the HTTP server on the specified port
// It sets up the routes, adds the request logger middleware, and starts listening for requests
func (s *Server) Start(port int) error {
	mux := s.Routes()

	addr := fmt.Sprintf(":%d", port)
	s.logger.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, s.requestLogger(mux))
}
