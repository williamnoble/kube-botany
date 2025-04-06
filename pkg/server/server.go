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

type Server struct {
	StaticDir    string
	TemplatesDir string
	templates    map[string]*template.Template

	logger    *slog.Logger
	startTime time.Time

	mu     sync.Mutex
	plants []*plant.Plant

	renderer *render.ASCIIRenderer
}

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
		StaticDir:    "cmd/api/static",
		TemplatesDir: "cmd/api/templates",
	}
	s.ParseTemplates()

	return s
}

func (s *Server) Start(port int) error {
	mux := s.Routes()

	addr := fmt.Sprintf(":%d", port)
	s.logger.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, s.requestLogger(mux))
}
