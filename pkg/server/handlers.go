package server

import (
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"github.com/williamnoble/kube-botany/pkg/types"
	"net/http"
	"time"
)

// HandleListPlants returns a list of all plants as JSON
func (s *Server) HandleListPlants(w http.ResponseWriter, r *http.Request) {
	var plants []types.PlantDTO
	for _, currentPlant := range s.store.ListAllPlants() {
		plantDTO := types.IntoPlantDTO(currentPlant)
		plants = append(plants, plantDTO)
	}
	err := s.encodeJsonResponse(w, r, http.StatusOK, plants)
	if err != nil {
		s.InternalServerErrorResponse(w, err)
	}
}

// HandleGetPlant returns a single plant by ID as JSON
func (s *Server) HandleGetPlant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := s.store.GetPlant(id)
	if err != nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}
	plantDTO := types.IntoPlantDTO(p)
	err = s.encodeJsonResponse(w, r, http.StatusOK, plantDTO)
	if err != nil {
		s.InternalServerErrorResponse(w, err)
	}
}

// HandleGetPlantAscii returns a single plant by ID as ASCII text including an ASCII image
// which varies depending on the plant's growth stage
func (s *Server) HandleGetPlantAscii(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := s.store.GetPlant(id)
	if err != nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	text := s.renderer.RenderText(p)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte(text))
	if err != nil {
		s.InternalServerErrorResponse(w, err)
	}
}

// HandlePlantDelete deletes a plant by ID
// It returns 204 No Content if successful, or 404 Not Found if the plant doesn't exist
func (s *Server) HandlePlantDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.store.DeletePlant(id)
	if err != nil {
		s.InternalServerErrorResponse(w, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// HandleWaterPlant adds water; a single request waters a plant to 100%
func (s *Server) HandleWaterPlant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := s.store.GetPlant(id)
	if err != nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	message := "plant is fully watered and cannot be watered anymore."
	unitsAdded := p.AddWater()
	if unitsAdded > 0 {
		message = fmt.Sprintf("added %d units of water to %s (%d%% watered).", unitsAdded, p.Id, p.CurrentWaterLevel())
	}

	response := WaterResponse{
		Message: message,
		Plant:   types.IntoPlantDTO(p),
	}

	err = s.encodeJsonResponse(w, r, http.StatusOK, response)
	if err != nil {
		s.InternalServerErrorResponse(w, err)
	}
}

// HandleRenderHomePage renders the home page with cards for all plants
// It converts each plant to a DTO, sets the image path, and renders the index.html template
func (s *Server) HandleRenderHomePage(w http.ResponseWriter, r *http.Request) {
	var data []types.PlantDTO
	plants := s.store.ListAllPlants()
	for _, plant := range plants {
		dto := types.IntoPlantDTO(plant)
		if dto.FriendlyName == "" {
			dto.FriendlyName = dto.Id
		}
		data = append(data, dto)
	}

	// Execute the index.html template with the layout
	err := s.templates["index"].ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.Logger.Error("template error", "error", err)
	}
}

// HandlePlantDetail renders the plant detail page for a specific plant
// It extracts the plant ID from the URL path, finds the plant, and renders the plant.html template
func (s *Server) HandlePlantDetail(w http.ResponseWriter, r *http.Request) {
	// Extract plant ID from the URL path
	id := r.PathValue("id")

	p, err := s.store.GetPlant(id)
	if err != nil {
		http.Error(w, "Plant not found", http.StatusNotFound)
		return
	}

	plantDTO := types.IntoPlantDTO(p)
	if plantDTO.FriendlyName == "" {
		plantDTO.FriendlyName = plantDTO.Id
	}

	err = s.templates["plant"].ExecuteTemplate(w, "layout.html", plantDTO)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		s.Logger.Error("template error", "error", err)
	}
}

// HandleCreatePlant creates a new plant from the request body
func (s *Server) HandleCreatePlant(w http.ResponseWriter, r *http.Request) {
	var dto types.PlantDTO
	err := s.decodeJsonRequest(r, &dto)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.store.NewPlant(
		dto.Id,
		dto.FriendlyName,
		dto.Variety,
		time.Now(),
	)

	if err != nil {
		http.Error(w, "Error creating plant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
