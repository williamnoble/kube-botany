package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	var response []PlantDTO
	for _, p := range s.plants {
		p.Update(time.Now())
		plant := s.plantDTO(p)
		response = append(response, plant)
	}

	err := s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandlePlant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var found bool
	var response PlantDTO
	for _, p := range s.plants {
		if p.Id == id {
			found = true
			p.Update(time.Now())
			response = s.plantDTO(p)
		}
	}

	if found == false || id == "" {
		http.Error(w, "Plant not found, please check the id.", http.StatusNotFound)
		return
	}

	err := s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandlePlantDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var found bool
	var index int
	for i, p := range s.plants {
		if p.Id == id {
			found = true
			p.Update(time.Now())
			index = i
			break
		}
	}

	if found == false || id == "" {
		http.Error(w, "Plant not found, please check the id.", http.StatusNotFound)
		return
	}

	s.plants = append(s.plants[:index], s.plants[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) HandleWater(w http.ResponseWriter, r *http.Request) {
	var waterReq WaterRequest
	err := s.decode(r, &waterReq)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	plant, err := s.plantByID(waterReq.Id)
	if err != nil {
		http.Error(w, "404", http.StatusBadRequest)
		return
	}

	message := "plant is fully watered and cannot be watered anymore."
	unitsAdded := plant.AddWater(time.Now())
	if unitsAdded > 0 {
		message = fmt.Sprintf("added %d units of water to %s (%d%% watered).", unitsAdded, plant.Id, plant.WaterLevel)
	}

	response := WaterResponse{
		Message: message,
		Plant:   s.plantDTO(plant),
	}

	err = s.encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) HandleCards(w http.ResponseWriter, r *http.Request) {
	var data []PlantDTO
	for _, plant := range s.plants {
		dto := s.plantDTO(plant)
		dto.Image = fmt.Sprintf("/static/images/%s", plant.Image())
		if dto.Name == "" {
			dto.Name = dto.Id
		}
		data = append(data, dto)
	}

	// Make sure you're executing the layout template
	err := s.templates["index"].ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.logger.Error("template error", "error", err)
	}
}

// HandlePlantDetail Adds a new handler for plant detail pages
func (s *Server) HandlePlantDetail(w http.ResponseWriter, r *http.Request) {
	// Extract plant name from URL path
	plantId := r.PathValue("id")

	var dto PlantDTO
	for _, plant := range s.plants {
		if plant.Id == plantId {
			dto = s.plantDTO(plant)
			dto.Image = fmt.Sprintf("/static/images/%s", plant.Image())
			if dto.Name == "" {
				dto.Name = dto.Id
			}
		}
	}

	// Execute only the plant.html template
	err := s.templates["plant"].ExecuteTemplate(w, "layout.html", dto)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.logger.Error("template error", "error", err)
	}
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
