package main

import (
	"encoding/json"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"github.com/williamnoble/kube-botany/pkg/render"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	s := NewServer(plant.Fern, "MyFern")
	if err := s.Start(8080); err != nil {
		panic(err)
	}
}

type Server struct {
	logger    *slog.Logger
	renderer  *render.ASCIIRenderer
	startTime time.Time
	plants    []*plant.Plant
}

type PlantResponse struct {
	ID          string    `json:"ID"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Age         string    `json:"age"`
	Growth      int64     `json:"growth"`
	GrowthStage string    `json:"stage"`
	WaterLevel  string    `json:"water"`
	WateredLast time.Time `json:"watered_last"`
}

type WaterResponse struct {
	Message string `json:"message"`
	Plant   PlantResponse
}

type WaterRequest struct {
	Id string `json:"id"` // NamespacedName
}

func NewServer(plantType plant.Type, plantName string) *Server {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)
	var plants []*plant.Plant
	plants = append(plants, plant.NewPlant(plantName, plantType, false))
	return &Server{
		plants:    plants,
		logger:    logger,
		startTime: time.Now(),
		renderer:  render.NewASCIIRenderer(),
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.HandleIndex)
	mux.HandleFunc("POST /water", s.HandleWater)
	//mux.HandleFunc("POST /plant", s.handleNewPlant)
	//mux.HandleFunc("GET /ascii", s.handleASCII)

	addr := fmt.Sprintf(":%d", port)
	s.logger.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, s.RequestLogger(mux))
}

func (s *Server) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		s.logger.Info("request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		s.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

//
//func (s *Server) handleASCII(w http.ResponseWriter, r *http.Request) {
//	currentPlant.Update(time.Now())
//	currentPlant.GrowthStage = plant.Maturing
//	asciiArt := s.renderer.RenderFern(currentPlant)
//	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(asciiArt))
//}

func (s *Server) encode(w http.ResponseWriter, r *http.Request, status int, v interface{}, wrap ...string) error {
	if len(wrap) > 0 {
		v = map[string]interface{}{
			wrap[0]: v,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func (s *Server) decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.logger.Error("failed to decode water request", "error", err)
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func (s *Server) plantByID(id string) (*plant.Plant, error) {
	for _, p := range s.plants {
		if p.ID == id {
			return p, nil
		}
	}

	return nil, fmt.Errorf("plant not found")
}

func (s *Server) plantResponse(p *plant.Plant) PlantResponse {
	r := PlantResponse{
		ID:          p.ID,
		Name:        p.Name,
		Type:        string(p.Type),
		GrowthStage: p.GrowthStage.String(),
		Age:         time.Since(p.CreationTime).Round(time.Second).String(),
		WaterLevel:  p.WaterLevelPercent(),
		WateredLast: p.LastWatered,
	}
	return r
}
