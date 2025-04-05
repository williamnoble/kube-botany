package main

import (
	"encoding/json"
	"fmt"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	s := NewServer(plant.Fern, "My Fern")
	if err := s.Start(8080); err != nil {
		panic(err)
	}
}

type Server struct {
	logger    *slog.Logger
	renderer  *ASCIIRenderer
	startTime time.Time
	plant     *plant.Plant
}

type PlantResponse struct {
	ID          string  `json:"ID"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Age         string  `json:"age"`
	GrowthStage string  `json:"growthStage"`
	Growth      float64 `json:"current_growth"`
	Health      string  `json:"health"`
	Water       struct {
		Min   float64 `json:"min" json:"min"`
		Max   float64 `json:"max" json:"max"`
		Level float64 `json:"level"`
	} `json:"water"`
}

type WaterResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	//Plant   PlantResponse `json:"plant"`
}

type WaterRequest struct {
	Id     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Amount float64 `json:"amount"`
}

func NewServer(plantType plant.Type, plantName string) *Server {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	return &Server{
		plant:     plant.NewPlant(plantName, plantType),
		logger:    logger,
		startTime: time.Now(),
		renderer:  NewASCIIRenderer(),
	}
}

func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("POST /water", s.handleWater)
	mux.HandleFunc("GET /ascii", s.handleASCII)

	addr := fmt.Sprintf(":%d", port)
	s.logger.Info("starting server", "addr", addr)
	return http.ListenAndServe(addr, s.requestLogger(mux))
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
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

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	s.plant.Update(time.Now())

	response := PlantResponse{
		ID:          s.plant.ID,
		Name:        s.plant.Name,
		Type:        string(s.plant.Type),
		GrowthStage: s.plant.GrowthStage.String(),
		Growth:      s.plant.Growth,
		Health:      s.plant.HealthPercent(),
		Age:         time.Since(s.plant.CreationTime).Round(time.Second).String(),
	}

	response.Water.Min = s.plant.Characteristics.OptimalWaterMin
	response.Water.Max = s.plant.Characteristics.OptimalWaterMax
	response.Water.Level = s.plant.WaterLevelFormatted()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleWater(w http.ResponseWriter, r *http.Request) {
	var waterReq WaterRequest

	if err := json.NewDecoder(r.Body).Decode(&waterReq); err != nil {
		s.logger.Error("failed to decode water request", "error", err)
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	plantID := waterReq.Id
	plantName := waterReq.Name

	if plantID != "" {
		s.logger.Info("plant identifier provided", "plant_id", plantID)
	} else if plantName != "" {
		s.logger.Info("plant name provided", "plant_name", plantName)
	}

	reqAmount := waterReq.Amount
	s.plant.Update(time.Now())
	actualAmount := s.plant.Water(reqAmount, time.Now())

	// Create the appropriate message based on how much was actually added
	var message string
	if actualAmount == 0 {
		// No water was added
		message = "Plant is already fully watered and cannot absorb any more water."
	} else if actualAmount == waterReq.Amount {
		// All water was added
		message = fmt.Sprintf("Successfully watered plant with %.2f units.", waterReq.Amount)
	} else {
		// Some water was added, but not all
		overflow := plant.RoundToTwoDecimal(waterReq.Amount - actualAmount)
		message = fmt.Sprintf("Plant absorbed %.2f units of water. %.2f units overflowed.",
			actualAmount, overflow)
	}

	s.logger.Info("watered plant",
		"plant", s.plant.Name,
		"reqAmount", reqAmount,
		"new_level", s.plant.WaterLevelFormatted(),
		"plant_health", s.plant.HealthPercent())

	plantResponse := PlantResponse{
		ID:          s.plant.ID,
		Name:        s.plant.Name,
		Health:      s.plant.HealthPercent(),
		GrowthStage: s.plant.GrowthStage.String(),
	}

	plantResponse.Water.Level = s.plant.WaterLevelFormatted()

	response := WaterResponse{
		Success: true,
		Message: message,
		//Plant:   plantResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleASCII(w http.ResponseWriter, r *http.Request) {

	// Update plant state to current time
	s.plant.Update(time.Now())

	// For now, we only support rendering ferns in ASCII
	asciiArt := s.renderer.RenderFern(s.plant)

	// Log the ASCII art for debugging
	s.logger.Info("rendering ASCII art",
		"plant_name", s.plant.Name,
		"growth_stage", s.plant.GrowthStage,
		"health", s.plant.Health,
		"water_level", s.plant.WaterLevel)

	// Set plain text content type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(asciiArt))
}
